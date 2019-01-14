package dbs

import "time"

const (
	//Deleted be soft deleted in dbs
	Deleted = 1
	//Undeleted inverse
	Undeleted = 0
)

type LibUnit struct {
	//ID desperated
	ID         int64  `xorm:"bigint(20) notnull pk autoincr 'id'"`
	Name       string `xorm:"varchar(256) notnull default '' 'name'"`
	RoomID     string `xorm:"varchar(256) notnull default '' 'room_id'"`
	Deleted    int8   `xorm:"tinyint(4) notnull default 0 'deleted'"`
	UpdateTime int64  `xorm:"bigint(20) notnull default 0 'update_time'"`
}

func (unit *LibUnit) TableName() string {
	return "lib_unit"
}

type User struct {
	//ID desperated
	ID       int64  `xorm:"pk autoincr comment('the table unique id') BIGINT(20) 'id'"`
	UserID   int64  `xorm:"not null default 0 comment('the user id specified from outer information') BIGINT(20) 'user_id'"`
	UserName string `xorm:"not null default '' VARCHAR(256) 'user_name'"`
	Deleted  int    `xorm:"not null default 0 TINYINT(4)"`
}

type Seat struct {
	//ID desperated
	ID       int64  `xorm:"pk autoincr BIGINT(20) 'id'"`
	SeatName string `xorm:"not null default '' VARCHAR(256)"`
	Location string `xorm:"not null default '' VARCHAR(256)"`
	RoomID   string `xorm:"not null default '' VARCHAR(256) 'room_id'"`
	Deleted  int    `xorm:"not null default 0 TINYINT(4)"`
}

type UserSeat struct {
	ID         int64     `xorm:"pk autoincr BIGINT(20) 'id'"`
	UserID     int64     `xorm:"not null default 0 BIGINT(20) 'user_id'"`
	SeatName   string    `xorm:"not null default '' VARCHAR(256)"`
	StartTime  time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	UpdateTime time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
	Deleted    int       `xorm:"not null default 0 TINYINT(4)"`
}
