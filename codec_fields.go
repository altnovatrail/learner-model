package learnermodel

// The field numbers of the schema's messages, proto/novatrail/learnermodel/v1/, each named after
// its message and its field there: the writers and the readers of the codec take a field by its
// name, and a number here changes only with the schema.

// BandRef.
const (
	fieldBandPlace = 1
	fieldBandID    = 2
)

// AnswerSlot.
const (
	fieldSlotCheckKind   = 1
	fieldSlotOptionCount = 2
)

// TrapCount.
const (
	fieldTrapCountID            = 1
	fieldTrapCountEvidenceCount = 2
	fieldTrapCountLastFoundTime = 3
)

// Evidence.
const (
	fieldEvidenceKind                     = 1
	fieldEvidenceTaskID                   = 2
	fieldEvidenceTemplate                 = 3
	fieldEvidenceTaskBand                 = 4
	fieldEvidenceDifficulty               = 5
	fieldEvidenceStartBand                = 6
	fieldEvidencePrimarySkills            = 7
	fieldEvidenceSupportingSkills         = 8
	fieldEvidenceKeyProvenance            = 9
	fieldEvidenceSkillTagsProvenance      = 10
	fieldEvidenceOutcome                  = 11
	fieldEvidenceHintRung                 = 12
	fieldEvidencePracticeTime             = 13
	fieldEvidenceTechniqueOnStep          = 14
	fieldEvidenceTechniqueOnRequest       = 15
	fieldEvidenceExplainAgainCount        = 16
	fieldEvidenceSolutionSeenBefore       = 17
	fieldEvidenceItems                    = 18
	fieldEvidenceTraps                    = 19
	fieldEvidenceLeadingQuestionOnRequest = 20
	fieldEvidenceLevelSetByAdult          = 21
)

// ItemEvidence.
const (
	fieldItemPlace                   = 1
	fieldItemOwnSkills               = 2
	fieldItemPrimarySkills           = 3
	fieldItemSupportingSkills        = 4
	fieldItemAnswerSlots             = 5
	fieldItemSubmitted               = 6
	fieldItemUnverdictedReasons      = 7
	fieldItemFirstVerdict            = 8
	fieldItemWrongBeforeFirstVerdict = 9
	fieldItemFinalVerdict            = 10
	fieldItemWrongBeforeFirstCorrect = 11
)

// FoundTrap.
const (
	fieldFoundTrapID         = 1
	fieldFoundTrapSkills     = 2
	fieldFoundTrapItemPlaces = 3
)

// State.
const (
	fieldStateKind    = 1
	fieldStateModules = 2
)

// ModuleState.
const (
	fieldModuleID                  = 1
	fieldModuleLevel               = 2
	fieldModuleAnswerWeight        = 3
	fieldModuleStartBand           = 4
	fieldModuleTrialAnswers        = 5
	fieldModuleTraps               = 6
	fieldModuleSkills              = 7
	fieldModuleLastEvidenceInTrial = 8
	fieldModuleFailuresInARow      = 9
	fieldModuleTrialEndedByAdult   = 10
)

// TrialAnswer.
const (
	fieldTrialAnswerBandPlace   = 1
	fieldTrialAnswerDifficulty  = 2
	fieldTrialAnswerAnswerSlots = 3
	fieldTrialAnswerVerdict     = 4
	fieldTrialAnswerWeight      = 5
)

// SkillState.
const (
	fieldSkillID               = 1
	fieldSkillCorrection       = 2
	fieldSkillAnswerWeight     = 3
	fieldSkillEvidenceCount    = 4
	fieldSkillLastPracticeTime = 5
	fieldSkillStability        = 6
	fieldSkillTraps            = 7
	fieldSkillMasteryBandPlace = 8
	fieldSkillFailuresInARow   = 9
)
