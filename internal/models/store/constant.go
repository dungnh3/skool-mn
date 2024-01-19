package store

import "encoding/json"

type CustomRawMessage json.RawMessage

type ObjectStatus string

const (
	Active   ObjectStatus = "active"
	InActive ObjectStatus = "inactive"
)

func (s ObjectStatus) String() string { return string(s) }

type RegisterStatus string

const (
	Registered       RegisterStatus = "registered"
	Confirmed        RegisterStatus = "confirmed"
	Cancelled        RegisterStatus = "cancelled"
	Rejected         RegisterStatus = "rejected"
	Done             RegisterStatus = "done"
	Waiting          RegisterStatus = "waiting"
	StudentLeftClass RegisterStatus = "student_left_class"
	StudentOutSchool RegisterStatus = "student_out_school"
)

func (s RegisterStatus) String() string { return string(s) }

type ActionType string

const (
	ActionRegis         ActionType = "regis"
	ActionConfirm       ActionType = "confirm"
	ActionCancel        ActionType = "cancel"
	ActionReject        ActionType = "reject"
	ActionWait          ActionType = "wait"
	ActionLeaveClass    ActionType = "leave_class"
	ActionOutSchool     ActionType = "out_school"
	ActionParentConfirm ActionType = "parent_confirm"
)

func (s ActionType) String() string { return string(s) }
