package model

// Consolidated enum set. Java carries two overlapping/inconsistent enum packages
// (web.models.* vs web.models.enums.*, with two conflicting Source enums) —
// the Go rewrite keeps a single authoritative set instead of replicating that split.
type Status string

const (
	StatusActive     Status = "ACTIVE"
	StatusInactive   Status = "INACTIVE"
	StatusInWorkflow Status = "INWORKFLOW"
)

type Channel string

const (
	ChannelSystem     Channel = "SYSTEM"
	ChannelCFCCounter Channel = "CFC_COUNTER"
	ChannelCitizen    Channel = "CITIZEN"
	ChannelDataEntry  Channel = "DATA_ENTRY"
	ChannelMigration  Channel = "MIGRATION"
)

type Relationship string

const (
	RelationshipFather  Relationship = "FATHER"
	RelationshipHusband Relationship = "HUSBAND"
)
