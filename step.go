package learnermodel

import (
	"cmp"
	"math"
	"slices"
	"strings"
)

// The step of one piece of evidence (docs/architecture.md, sections 3.1, 3.3 and 3.5): a function
// of the module's state before it, the evidence and the graph of the module's current version,
// the same when appending and when recomputing. Every P and K comes from the state before the
// evidence, and the moves of its answers add up into one step. Sums run in the canonical order —
// items by place, skills by ID — so that the bits do not depend on the order of the input.
//
// From T16 the step is that of R03 with the weights of R06: the first verdict of an item is its
// answer, each of an attempt's N answers weighs u = 3 / (N + 2), a primary skill of an item gets
// 1/J of its step, and a failure is fewer than half of a skill's answers right. Retention and
// stability (R04), help, outcomes, provenance and traps (R06) and the transfer by retired_skills
// (R08) come with T17–T19: until then R̄ = 1, every right answer is unaided, and a primary skill
// that is not a skill of the graph with a verdict gets nothing.

// item is an item of the task with what happened to it in the attempt.
type item struct {
	task   TaskItem
	answer ItemAnswer
}

// pairItems pairs the items of the evidence's task with its answers, by place. The evidence has
// passed checkEvidenceContent, so the two name the same places, each once.
func pairItems(e Evidence) []item {
	tasks := slices.Clone(e.Task.Items)
	slices.SortFunc(tasks, func(a, b TaskItem) int { return cmp.Compare(a.Place, b.Place) })
	answers := slices.Clone(e.Items)
	slices.SortFunc(answers, func(a, b ItemAnswer) int { return cmp.Compare(a.Place, b.Place) })
	items := make([]item, 0, len(tasks))
	for i, t := range tasks {
		items = append(items, item{task: t, answer: answers[i]})
	}
	return items
}

// answer is an item of the evidence that gives an answer, as the step reads it.
type answer struct {
	slots   []Slot
	verdict Verdict
	// skills are the item's primary skills that have a state to move, by ID: J = len(skills).
	skills []Skill
	p      float64
}

func (a answer) outcome() float64 {
	if a.verdict == VerdictCorrect {
		return 1
	}
	return 0
}

// skillMove gathers what one piece of evidence does to a skill.
type skillMove struct {
	skill   Skill
	answers int     // N_s, the answers where it is a primary skill
	correct int     // the right ones among them
	pull    float64 // Σ (S − P) / J over them
}

// itemSkills gives the primary skills of an item: its own, or the task's.
func itemSkills(t Task, item TaskItem) []string {
	if item.OwnSkills {
		return item.Primary
	}
	return t.Primary
}

// isAnswer tells whether an item gives an answer: a first verdict right or not yet, and every
// slot of a kind this version knows (docs/architecture.md, section 3.5).
func isAnswer(item TaskItem, a ItemAnswer) bool {
	if !a.HasFirstVerdict || (a.FirstVerdict != VerdictCorrect && a.FirstVerdict != VerdictNotYet) {
		return false
	}
	for _, s := range item.Slots {
		if s.Check < CheckKindNumber || s.Check > CheckKindPredicate {
			return false
		}
	}
	return true
}

// graphSkill finds a skill of the graph by its ID.
func graphSkill(g *Graph, id string) (Skill, bool) {
	i, ok := slices.BinarySearchFunc(g.Skills, id, func(s Skill, id string) int { return strings.Compare(s.ID, id) })
	if !ok {
		return Skill{}, false
	}
	return g.Skills[i], true
}

// skillIndex finds a skill's state in a module, or where it would stand.
func (ms *ModuleState) skillIndex(id string) (int, bool) {
	return slices.BinarySearchFunc(ms.Skills, id, func(s SkillState, id string) int { return strings.Compare(s.Skill, id) })
}

func (ms *ModuleState) correction(id string) float64 {
	if i, ok := ms.skillIndex(id); ok {
		return ms.Skills[i].Correction
	}
	return 0
}

// answers reads the answers of the evidence, items by place, each with its P from the state
// before it: P = c + (1 − c) · σ(θ + δ̄ − β), δ̄ the mean δ of its J skills, 0 with none.
func (ms *ModuleState) answers(t Task, items []item, g *Graph) []answer {
	b := beta(t.Band.Place, t.Difficulty)
	var out []answer
	for _, it := range items {
		if !isAnswer(it.task, it.answer) {
			continue
		}
		var skills []Skill
		for _, id := range itemSkills(t, it.task) {
			if s, ok := graphSkill(g, id); ok && s.Assessability == AssessabilityVerdict {
				skills = append(skills, s)
			}
		}
		slices.SortFunc(skills, func(a, b Skill) int { return strings.Compare(a.ID, b.ID) })
		var mean float64
		if len(skills) > 0 {
			var sum float64
			for _, s := range skills {
				sum += ms.correction(s.ID)
			}
			mean = float64(sum / float64(len(skills)))
		}
		out = append(out, answer{
			slots:   it.task.Slots,
			verdict: it.answer.FirstVerdict,
			skills:  skills,
			p:       probability(guess(it.task.Slots), ms.Level+mean-b),
		})
	}
	return out
}

// moves gathers the evidence's answers by skill, in the order of skill IDs.
func moves(answers []answer) []skillMove {
	var out []skillMove
	for _, a := range answers {
		j := float64(len(a.skills))
		for _, s := range a.skills {
			i, ok := slices.BinarySearchFunc(out, s.ID, func(m skillMove, id string) int { return strings.Compare(m.skill.ID, id) })
			if !ok {
				out = slices.Insert(out, i, skillMove{skill: s})
			}
			m := &out[i]
			m.answers++
			if a.verdict == VerdictCorrect {
				m.correct++
			}
			m.pull += float64((a.outcome() - a.p) / j)
		}
	}
	return out
}

// step moves the module's state by one piece of evidence, against the graph of the module's
// current version: the trial series while the weight of answers is below 5, the step of Elo after
// it. Evidence with no answer moves nothing, and step reports whether it had one.
func (ms *ModuleState) step(e Evidence, items []item, g *Graph) bool {
	return ms.move(e, items, g, ms.AnswerWeight < trialWeight-weightTolerance)
}

// move is the step, with the trial series or without it as trial says.
func (ms *ModuleState) move(e Evidence, items []item, g *Graph, trial bool) bool {
	answers := ms.answers(e.Task, items, g)
	if len(answers) == 0 {
		return false
	}
	u := seriesWeight(len(answers))
	if trial {
		for _, a := range answers {
			ms.Trial = append(ms.Trial, TrialAnswer{
				BandPlace:  e.Task.Band.Place,
				Difficulty: e.Task.Difficulty,
				Slots:      slices.Clone(a.slots),
				Verdict:    a.verdict,
				Weight:     u,
			})
		}
		ms.Level = trialLevel(shift(ms.StartBand.Place), ms.Trial)
	} else {
		var pull float64
		for _, a := range answers {
			pull += a.outcome() - a.p
		}
		ms.Level += float64(levelStep(ms.AnswerWeight) * u * pull)
	}
	ms.AnswerWeight += float64(float64(len(answers)) * u)
	if trial && ms.AnswerWeight >= trialWeight-weightTolerance {
		ms.Trial = nil
	}

	for _, m := range moves(answers) {
		i, ok := ms.skillIndex(m.skill.ID)
		if !ok {
			ms.Skills = slices.Insert(ms.Skills, i, SkillState{Skill: m.skill.ID, LastPractice: e.PracticeTime})
		}
		sk := &ms.Skills[i]
		us := seriesWeight(m.answers)
		if !trial {
			sk.Correction += float64(correctionStep(sk.AnswerWeight) * us * m.pull)
		}
		sk.AnswerWeight += float64(float64(m.answers) * us)
		sk.EvidenceCount++
		sk.LastPractice = max(sk.LastPractice, e.PracticeTime)
		ms.judge(sk, m, e.Task.Band.Place)
	}
	return true
}

// judge counts a failure or sets the mastery after the evidence, by θ, δ, n and m after it. A
// failure takes the mastery away at the second in a row; evidence without one — which has a right
// answer, since at least half are — masters the skill once m ≥ 5, at a band above the one it is
// held at.
func (ms *ModuleState) judge(sk *SkillState, m skillMove, taskPlace int32) {
	if 2*m.correct < m.answers {
		if sk.FailuresInARow < math.MaxInt32 {
			sk.FailuresInARow++
		}
		if sk.FailuresInARow >= failuresToLose {
			sk.Mastered, sk.MasteryBand = false, 0
		}
		return
	}
	sk.FailuresInARow = 0
	if sk.AnswerWeight < masteryWeight-weightTolerance {
		return
	}
	lambda := cautiousLevel(ms.Level, sk.Correction, ms.AnswerWeight, sk.AnswerWeight)
	if p, ok := masteryBand(m.skill.BandsTaught, taskPlace, lambda); ok && (!sk.Mastered || p > sk.MasteryBand) {
		sk.Mastered, sk.MasteryBand = true, p
	}
}
