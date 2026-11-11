package learnermodel

import (
	"cmp"
	"math"
	"slices"
	"strings"

	"github.com/altnovatrail/learner-model/detmath"
)

// The step of one piece of evidence (docs/architecture.md, sections 3.1, 3.3 and 3.5): a function
// of the module's state before it, the evidence and the graph of the module's current version,
// the same when appending and when recomputing. Every P and K comes from the state before the
// evidence, and the moves of its answers add up into one step. Sums run in the canonical order —
// items by place, skills by ID — so that the bits do not depend on the order of the input.
//
// The step is that of R03 with the weights of R06 (T16): the first verdict of an item is its
// answer, each of an attempt's N answers weighs u = 3 / (N + 2), a primary skill of an item gets
// 1/J of its step, and a failure is fewer than half of a skill's answers right. Forgetting (R04,
// T17): P carries the mean retention R̄ of the item's skills, each answer moves θ and δ with its
// weight w, and a skill's stability τ changes once per piece of evidence, by the share of its
// answers right. The evidence read whole (R06, T18): the outcome, where a shown solution makes an
// empty item "not solved"; help, which masters nothing; the provenance of the key and the skill
// tags, by which evidence weighs 1 or nothing; and the traps the program found, counted at their
// skills or at the module. The transfer by retired_skills (R08, T19): each recorded primary skill
// of an item, and each skill of a trap, goes to the skills transfer gives it by the graph of the
// module's current version, so that a retired skill's evidence counts in full for each skill that
// replaced it, and skills merged in one item count once.

// item is an item of the task with what happened to it in the attempt, and the answer it gives by
// answerOf, worked out once for the checks of input and the step alike.
type item struct {
	task     TaskItem
	answer   ItemAnswer
	verdict  Verdict
	answered bool
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
		it := item{task: t, answer: answers[i]}
		it.verdict, it.answered = answerOf(e, it)
		items = append(items, it)
	}
	return items
}

// answer is an item of the evidence that gives an answer, as the step reads it.
type answer struct {
	slots   []Slot
	verdict Verdict
	// shares are the skills the item's primary skills go to after the transfer, by ID, each once.
	shares []share
	// j is J, the number of the item's recorded primary skills that go to some skill.
	j      int
	p      float64 // P, with the retention of the skills
	w      float64 // the weight of retention, w; 1 without forgetting
	helped bool    // given with help: it masters nothing
}

// share is a skill an answer moves, with k, the number of the item's recorded primary skills that
// go to it: it gets k/J of the item's step, and k is 2 for a skill that two of them merged into.
type share struct {
	skill Skill
	k     int
}

func (a answer) outcome() float64 {
	if a.verdict == VerdictCorrect {
		return 1
	}
	return 0
}

// skillMove gathers what one piece of evidence does to a skill.
type skillMove struct {
	skill     Skill
	retention float64 // R of the skill at the evidence's practice time, from its state before it
	answers   int     // N_s, the answers where it is a primary skill
	correct   int     // the right ones among them
	unaided   int     // the right ones among them given without help
	pull      float64 // Σ w · (S − P) · k / J over them
}

// itemSkills gives the primary skills of an item: its own, or the task's.
func itemSkills(t Task, item TaskItem) []string {
	if item.OwnSkills {
		return item.Primary
	}
	return t.Primary
}

// transfer gives the skills of the graph a recorded skill goes to (docs/architecture.md, section
// 3.5), by ID: the skill itself, when it is a skill of the graph with a verdict; the replacements
// with a verdict of a retired skill, which the graph names among its own skills, so that a chain
// of retirements is one step already; and none for a skill that now has a rubric or none, and for
// one the graph does not know — its evidence waits until a version names it. The skill itself is
// the graph's own element, with no copy, so the result is only read.
func transfer(g *Graph, id string) []Skill {
	if i, ok := graphSkillIndex(g, id); ok {
		if g.Skills[i].Assessability == AssessabilityVerdict {
			return g.Skills[i : i+1 : i+1]
		}
		return nil
	}
	i, ok := slices.BinarySearchFunc(g.Retired, id, func(r RetiredSkill, id string) int { return strings.Compare(r.ID, id) })
	if !ok {
		return nil
	}
	var out []Skill
	for _, r := range g.Retired[i].ReplacedBy {
		if s, ok := graphSkill(g, r); ok && s.Assessability == AssessabilityVerdict {
			out = append(out, s)
		}
	}
	return out
}

// knownChecks tells whether every slot is of a kind of check this version knows: an item with
// another gives no answer (docs/architecture.md, section 3.5).
func knownChecks(slots []Slot) bool {
	for _, s := range slots {
		if s.Check < CheckKindNumber || s.Check > CheckKindTime {
			return false
		}
	}
	return true
}

// answerOf gives the answer an item gives, by the outcome of the attempt (docs/architecture.md,
// section 3.3): its first verdict, right or not yet; or, when the solution was shown, "not
// solved" — not yet, S = 0 — for an item without one that was not submitted or only blank. An
// item that had only another "unreadable" or "unchecked" gives none: neither is a mistake. Nor
// does an item with a verdict or a kind of check this version does not know, and an unknown
// outcome reads as "unchecked": the verdicts before it count, and it weighs nothing.
func answerOf(e Evidence, it item) (Verdict, bool) {
	if !knownChecks(it.task.Slots) {
		return VerdictUnspecified, false
	}
	a := it.answer
	if a.HasFirstVerdict {
		if a.FirstVerdict == VerdictCorrect || a.FirstVerdict == VerdictNotYet {
			return a.FirstVerdict, true
		}
		return VerdictUnspecified, false
	}
	if e.Outcome.Kind != OutcomeSolutionShown {
		return VerdictUnspecified, false
	}
	if !a.Submitted {
		return VerdictNotYet, true
	}
	if len(a.Unverdicted) == 0 {
		return VerdictUnspecified, false
	}
	for _, r := range a.Unverdicted {
		if r != UnverdictedUnreadableBlank {
			return VerdictUnspecified, false
		}
	}
	return VerdictNotYet, true
}

// helped tells whether an item's answer was given with help (docs/architecture.md, section 3.3):
// a rung of the ladder before its first verdict — a wrong answer earlier in the attempt —, the
// technique or the leading question on request anywhere in the attempt, the solution of this task
// sent to the device before, or a sketch. The technique on the step and "explain again" are no
// help. An answer with help moves the step, the counts and τ as it is, but masters nothing.
func helped(e Evidence, it item) bool {
	// The codec keeps a count of wrong answers only beside a first verdict, and the solution seen
	// only beside its fact.
	return it.answer.WrongBeforeFirst > 0 || e.TechniqueOnRequest || e.LeadingQuestionOnRequest ||
		e.SolutionSeenBefore || e.Outcome.Kind == OutcomeWithSketch
}

// weighs tells whether the evidence weighs anything by where its key and skill tags came from
// (docs/architecture.md, section 3.3; concepts 9.4 of spec): both from the author, the solver or an
// adult — not from AI. A key or tags from AI weigh 0, and so does a provenance not given or one
// this version does not know: such evidence moves nothing, as evidence with no answer does. An
// adult's key needs a band and a difficulty on the module's scale too, and evidence always has
// them.
func weighs(t Task) bool {
	return trusted(t.KeyProvenance) && trusted(t.SkillTagsProvenance)
}

// trusted tells whether a provenance is one that weighs: the author, the solver or an adult.
func trusted(p Provenance) bool {
	return p == ProvenanceAuthor || p == ProvenanceSolver || p == ProvenanceAdult
}

// graphSkill finds a skill of the graph by its ID.
func graphSkill(g *Graph, id string) (Skill, bool) {
	i, ok := graphSkillIndex(g, id)
	if !ok {
		return Skill{}, false
	}
	return g.Skills[i], true
}

// graphSkillIndex finds where a skill of the graph stands, by its ID.
func graphSkillIndex(g *Graph, id string) (int, bool) {
	return slices.BinarySearchFunc(g.Skills, id, func(s Skill, id string) int { return strings.Compare(s.ID, id) })
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

// skillRetention is R of a skill at the evidence's practice time, from the state before it.
type skillRetention struct {
	skill string
	r     float64
}

// retention gives R of a skill at a time from the state before the evidence, worked out once a
// skill: rs keeps those already worked out, in the order of IDs, and gains this one.
func (ms *ModuleState) retention(rs *[]skillRetention, id string, at int64) float64 {
	i, found := slices.BinarySearchFunc(*rs, id, func(r skillRetention, id string) int { return strings.Compare(r.skill, id) })
	if !found {
		// A skill with no practice forgets nothing.
		r := 1.0
		if k, ok := ms.skillIndex(id); ok {
			r = retentionAt(ms.Skills[k], at)
		}
		*rs = slices.Insert(*rs, i, skillRetention{skill: id, r: r})
	}
	return (*rs)[i].r
}

// retentionOf gives R of a skill that answers has worked out.
func retentionOf(rs []skillRetention, id string) float64 {
	i, _ := slices.BinarySearchFunc(rs, id, func(r skillRetention, id string) int { return strings.Compare(r.skill, id) })
	return rs[i].r
}

// itemChance is what P of an item of a task comes of, from the module's state at a time — the
// practice time of evidence, or the time of reading.
type itemChance struct {
	j     int     // J
	c     float64 // the guessing
	sigma float64 // σ(θ + δ̄ − β)
	rbar  float64 // R̄
	p1, p float64 // P at R̄ = 1, and P
}

// chanceOf works out P of an item from the state: P = c + (1 − c) · R̄ · σ(θ + δ̄ − β), δ̄ and R̄
// the mean δ and R of its J recorded skills, in the order of their IDs — 0 and 1 with none — where
// δ and R of a recorded skill are the means over the skills it goes to, and R of each skill at the
// time is worked out once into rs. With no transfer each of those means is of one value, and
// exact, so the bits are those of a step without the transfer. In the trial series P is without
// retention, R̄ = 1. The step and Readout.Task both read P by it, so that reading gives the P the
// step would take of evidence at the time of reading. The step gives shares too, and gets the
// skills the answer moves into them; reading gives nil, and no share is gathered.
func (ms *ModuleState) chanceOf(t Task, ti TaskItem, g *Graph, rs *[]skillRetention, shares *[]share, at int64, trial bool) itemChance {
	var out itemChance
	var corrections, retentions float64
	for _, x := range sortedView(itemSkills(t, ti)) {
		to := transfer(g, x)
		if len(to) == 0 {
			continue
		}
		out.j++
		var correction, retention float64
		for _, s := range to {
			correction += ms.correction(s.ID)
			if r := ms.retention(rs, s.ID, at); !trial {
				retention += r
			}
			if shares != nil {
				i, ok := slices.BinarySearchFunc(*shares, s.ID, func(sh share, id string) int { return strings.Compare(sh.skill.ID, id) })
				if !ok {
					*shares = slices.Insert(*shares, i, share{skill: s})
				}
				(*shares)[i].k++
			}
		}
		parts := float64(len(to))
		corrections += float64(correction / parts)
		if !trial {
			retentions += float64(retention / parts)
		}
	}
	mean := 0.0
	out.rbar = 1
	if out.j > 0 {
		mean = float64(corrections / float64(out.j))
		if !trial {
			out.rbar = float64(retentions / float64(out.j))
		}
	}
	out.c = guess(ti.Slots)
	out.sigma = detmath.Sigmoid(ms.Level + mean - beta(t.Band.Place, t.Difficulty))
	out.p1, out.p = chance(out.c, out.sigma, 1), chance(out.c, out.sigma, out.rbar)
	return out
}

// answers reads the answers of the evidence, items by place, each with its P from the state
// before it, by chanceOf, and its weight w; in the trial series w = 1. It gives with them R of
// every skill they move, for the stability, in the order of IDs: R of each skill is the same for
// every item that names it.
func (ms *ModuleState) answers(e Evidence, items []item, g *Graph, trial bool) ([]answer, []skillRetention) {
	var out []answer
	var rs []skillRetention
	for _, it := range items {
		if !it.answered {
			continue
		}
		var shares []share
		ic := ms.chanceOf(e.Task, it.task, g, &rs, &shares, e.PracticeTime, trial)
		out = append(out, answer{
			slots:   it.task.Slots,
			verdict: it.verdict,
			helped:  helped(e, it),
			shares:  shares,
			j:       ic.j,
			p:       ic.p,
			w:       retentionWeight(ic.c, ic.sigma, ic.rbar, ic.p1, ic.p),
		})
	}
	return out, rs
}

// moves gathers the evidence's answers by skill after the transfer, in the order of skill IDs,
// each with its retention from rs: a skill gets k/J of an answer's step, and counts the answer
// once however many of the item's skills went to it.
func moves(answers []answer, rs []skillRetention) []skillMove {
	var out []skillMove
	for _, a := range answers {
		j := float64(a.j)
		for _, sh := range a.shares {
			s := sh.skill
			i, ok := slices.BinarySearchFunc(out, s.ID, func(m skillMove, id string) int { return strings.Compare(m.skill.ID, id) })
			if !ok {
				out = slices.Insert(out, i, skillMove{skill: s, retention: retentionOf(rs, s.ID)})
			}
			m := &out[i]
			m.answers++
			if a.verdict == VerdictCorrect {
				m.correct++
				if !a.helped {
					m.unaided++
				}
			}
			// k = 1 multiplies exactly, so without a merge the bits are those of w · (S − P) / J.
			m.pull += float64(a.w * (a.outcome() - a.p) * float64(sh.k) / j)
		}
	}
	return out
}

// step moves the module's state by one piece of evidence, against the graph of the module's
// current version: the trial series while the weight of answers is below 5, the step of Elo after
// it. Evidence with no answer moves nothing, nor does evidence that weighs nothing by its
// provenance, and step reports whether it moved anything.
func (ms *ModuleState) step(e Evidence, items []item, g *Graph) bool {
	return ms.move(e, items, g, ms.inTrial())
}

// inTrial tells whether the module is in its trial series: the weight of its answers is below 5,
// within ε (docs/architecture.md, section 3.5).
func (ms *ModuleState) inTrial() bool {
	return ms.AnswerWeight < trialWeight-weightTolerance
}

// move is the step, with the trial series or without it as trial says.
func (ms *ModuleState) move(e Evidence, items []item, g *Graph, trial bool) bool {
	if !weighs(e.Task) {
		return false
	}
	at := e.PracticeTime
	answers, rs := ms.answers(e, items, g, trial)
	if len(answers) == 0 {
		return false
	}
	// The retention of each skill comes from the state before the evidence, as P does.
	skills := moves(answers, rs)
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
			pull += float64(a.w * (a.outcome() - a.p))
		}
		ms.Level += float64(levelStep(ms.AnswerWeight) * u * pull)
	}
	ms.AnswerWeight += float64(float64(len(answers)) * u)
	ms.LastInTrial = trial
	if trial && !ms.inTrial() {
		ms.Trial = nil
	}

	for _, m := range skills {
		i, ok := ms.skillIndex(m.skill.ID)
		if ok {
			sk := &ms.Skills[i]
			sk.Stability = stabilityAfter(sk.Stability, m.retention, float64(float64(m.correct)/float64(m.answers)))
		} else {
			// The first answer of a skill: τ = τ₀, and the time of practice starts.
			ms.Skills = slices.Insert(ms.Skills, i, SkillState{Skill: m.skill.ID, LastPractice: at, Stability: stabilityFloor})
		}
		sk := &ms.Skills[i]
		us := seriesWeight(m.answers)
		if !trial {
			sk.Correction += float64(correctionStep(sk.AnswerWeight) * us * m.pull)
		}
		sk.AnswerWeight += float64(float64(m.answers) * us)
		sk.EvidenceCount++
		// A time that came back winds no clock back.
		sk.LastPractice = max(sk.LastPractice, at)
		ms.judge(sk, m, e.Task.Band.Place)
	}
	ms.countTraps(e.Traps, g, at)
	return true
}

// judge counts a failure or sets the mastery after the evidence, by θ, δ, n and m after it. A
// failure takes the mastery away at the second in a row; evidence without one, where the skill has
// a right answer given without help, masters it once m ≥ 5, at a band above the one it is held at.
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
	if m.unaided == 0 || sk.AnswerWeight < masteryWeight-weightTolerance {
		return
	}
	lambda := cautiousLevel(ms.Level, sk.Correction, ms.AnswerWeight, sk.AnswerWeight)
	if p, ok := masteryBand(m.skill.BandsTaught, taskPlace, lambda); ok && (!sk.Mastered || p > sk.MasteryBand) {
		sk.Mastered, sk.MasteryBand = true, p
	}
}

// countTraps counts the traps the program found in the evidence, once each per piece of evidence
// however many items showed it (docs/architecture.md, sections 3.3 and 3.5): its skills go to
// skills by transfer, as an item's do, and the trap goes to those of them that have a state after
// the step, each once, with the number of pieces of evidence it was found in and the practice time
// of the latest; with no such skill, to the module. A trap gives no skill a state of its own: a
// skill has one only once it is practised.
func (ms *ModuleState) countTraps(traps []FoundTrap, g *Graph, at int64) {
	for _, trap := range traps {
		var states []int // the indexes of the skills' states, ascending, each once
		for _, x := range trap.Skills {
			for _, s := range transfer(g, x) {
				if i, ok := ms.skillIndex(s.ID); ok {
					if k, found := slices.BinarySearch(states, i); !found {
						states = slices.Insert(states, k, i)
					}
				}
			}
		}
		if len(states) == 0 {
			ms.Traps = countTrap(ms.Traps, trap.ID, at)
		}
		for _, i := range states {
			ms.Skills[i].Traps = countTrap(ms.Skills[i].Traps, trap.ID, at)
		}
	}
}

// countTrap counts one more piece of evidence of a trap in a list ordered by ID: a time that came
// back winds the latest back no more than it does a skill's practice.
func countTrap(traps []TrapCount, id string, at int64) []TrapCount {
	i, ok := slices.BinarySearchFunc(traps, id, func(t TrapCount, id string) int { return strings.Compare(t.ID, id) })
	if !ok {
		return slices.Insert(traps, i, TrapCount{ID: id, Count: 1, Last: at})
	}
	if traps[i].Count < math.MaxInt64 {
		traps[i].Count++
	}
	traps[i].Last = max(traps[i].Last, at)
	return traps
}
