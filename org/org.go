package org

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Org struct {
	ID               primitive.ObjectID   `bson:"_id,omitempty"`
	Slug             string               `bson:"s,omitempty"`
	DaDataID         string               `bson:"d,omitempty"`
	AreaID           primitive.ObjectID   `bson:"a,omitempty"`
	LocationID       primitive.ObjectID   `bson:"l,omitempty"`
	ManagerID        primitive.ObjectID   `bson:"mi,omitempty"`
	ManagerPost      string               `bson:"mp,omitempty"`
	EmployeeCount    uint32               `bson:"e,omitempty"`
	OkvedOsnID       primitive.ObjectID   `bson:"o,omitempty"`
	OkvedDopIDs      []primitive.ObjectID `bson:"od,omitempty"`
	Metros           []Metro              `bson:"m,omitempty"`
	Name             string               `bson:"n,omitempty"`
	NameFullWithOPF  string               `bson:"nf,omitempty"`
	NameShortWithOPF string               `bson:"ns,omitempty"`
	OPFCode          uint64               `bson:"oc,omitempty"`
	OPFFull          string               `bson:"of,omitempty"`
	OPFShort         string               `bson:"os,omitempty"`
	OPFKind          opfKind              `bson:"opk,omitempty"`
	Kind             kind                 `bson:"k,omitempty"`
	BranchKind       branchKind           `bson:"bk,omitempty"`
	BranchCount      uint32               `bson:"bc,omitempty"`
	INN              uint64               `bson:"i,omitempty"`
	KPP              uint64               `bson:"kp,omitempty"`
	OGRN             uint64               `bson:"og,omitempty"`
	OGRNDate         time.Time            `bson:"ogd,omitempty"`
	OKATO            uint64               `bson:"oka,omitempty"`
	OKTMO            uint64               `bson:"okt,omitempty"`
	OKPO             uint64               `bson:"okp,omitempty"`
	OKOGU            uint64               `bson:"oko,omitempty"`
	OKFS             uint64               `bson:"okf,omitempty"`
	StatusKind       statusKind           `bson:"sk,omitempty"`
	RegistrationDate time.Time            `bson:"rd,omitempty"`
	LiquidationDate  time.Time            `bson:"ld,omitempty"`
	UpdatedAt        time.Time            `bson:"ua,omitempty"`
}

type Metro struct {
	ID       primitive.ObjectID `bson:"id,omitempty"`
	Distance float32            `bson:"d,omitempty"`
}

type kind uint8

const (
	_ kind = iota
	kind_legal
	kind_individual
)

type branchKind uint8

const (
	_ branchKind = iota
	branchKind_main
	branchKind_branch
)

type opfKind uint8

const (
	_ opfKind = iota
	opfKind_y1999
	opfKind_y2012
	opfKind_y2014
)

type statusKind uint8

const (
	_ statusKind = iota
	statusKind_active
	statusKind_liquidating
	statusKind_liquidated
	statusKind_bankrupt
	statusKind_reorganizing
)
