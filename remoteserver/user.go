package remoteserver

type UserStatus uint8

const (
	// UserStatusPrepare 预录入数据
	UserStatusPrepare UserStatus = iota + 1
	// UserStatusNoRegistry 未注册
	UserStatusNoRegistry
	// UserStatusNormal 正常
	UserStatusNormal
	// UserStatusStop 停用
	UserStatusStop
	// UserStatusNotAllow 禁用
	UserStatusNotAllow
)

type WorkerStatus uint8

const (
	// WorkerStatusJon 在职
	WorkerStatusJon WorkerStatus = iota + 1
	// WorkerStatusLevel 离职
	WorkerStatusLevel
)

// Sex 性别
type Sex uint8

const (
	// SexMan 男
	SexMan Sex = iota + 1
	// SexWoman 女
	SexWoman
)

// UserDeptMain 部门主管信息
type UserDeptMain struct {
	// UserId 用户Id
	UserId string `json:"userId,omitempty" gorm:"primary_key"`
	// DeptId 部门ID
	DeptId string `json:"deptId,omitempty" gorm:"primary_key"`
	// Department 部门信息
	Department *Department `json:"department,omitempty" gorm:"foreignKey:deptId"`
	// Label 标签
	Label string `json:"name,omitempty"`
}

// Post 岗位表
type Post struct {
	// DbCommonField 通用字段
	Id string `json:"id,omitempty" gorm:"primary_key"`
	// Name 部门名称
	Name string `json:"name,omitempty" gorm:"not null"`
	// DeptId 部门ID
	DeptId string `json:"dept_id,omitempty"`
	// Department 部门信息
	Department *Department `gorm:"foreignKey:deptId"`
}

// Jobs 职位表
type Jobs struct {
	// Id 职位ID
	Id string `json:"id,omitempty" gorm:"primary_key"`
	// Name 部门名称
	Name string `json:"name,omitempty" gorm:"not null"`
	// DeptId 部门ID
	DeptId string `json:"dept_id,omitempty"`
	// Department 部门信息
	Department *Department `gorm:"foreignKey:deptId"`
}

// Department 部门表
type Department struct {
	// Id 部门ID
	Id string `json:"id,omitempty" gorm:"primary_key"`
	// Name 部门名称
	Name string `json:"name,omitempty" gorm:"not null;unique"`
	// Pid 上级部门ID
	Pid string `json:"pid,omitempty"`
	// Jobs 职位
	Jobs []*Jobs `json:"jobList,omitempty" gorm:"foreignKey:dept_id"`
	// Posts 岗位列表
	Posts []*Post `json:"postList,omitempty" gorm:"foreignKey:dept_id"`
	// Users 人员列表
	Users []*UserInfo `json:"userList,omitempty" gorm:"foreignKey:user_dept"`
	// MainList 部门主管人员列表
	MainList []*UserDeptMain `json:"MainList,omitempty" gorm:"foreignKey:dept_id"`
}

// UserDeptInfo 用户所在部门信息
type UserDeptInfo struct {
	// UserId 用户ID
	UserId string `json:"userId,omitempty" gorm:"primary_key"`
	// UserInfo 用户信息
	UserInfo *UserInfo `json:"userInfo,omitempty" gorm:"foreignKey:userId"`
	// DeptId 部门Id
	DeptId string `json:"deptId,omitempty" gorm:"primary_key;"`
	// Department 部门信息
	Department *Department `json:"department,omitempty" gorm:"foreignKey:deptId"`
	// Jobs 职位列表
	Jobs []*Jobs `json:"jobs,omitempty" gorm:"many2many:db_user_dept_jobs;foreignKey:deptId;"`
	// Posts 岗位列表
	Posts []*Post `json:"posts,omitempty" gorm:"many2many:db_user_dept_posts;foreignKey:deptId"`
}

// Staff 员工类型
type Staff struct {
	// Id 员工类型ID
	Id string `json:"Id,omitempty" gorm:"primary_key"`
	// Name 员工类型名称
	Name string `json:"name,omitempty" gorm:"not null"`
}

// Nation 名族
type Nation struct {
	// Id 名族类型ID
	Id string `json:"Id,omitempty" gorm:"primary_key"`
	// Name 名族名称
	Name string `json:"name,omitempty" gorm:"not null"`
}

// UserInfo 用户信息
type UserInfo struct {
	// Id 用户id
	Id string `json:"id,omitempty" gorm:"primary_key"`
	// Name 人员姓名
	Name string `json:"name,omitempty" gorm:"not null"`
	// Username 用户名
	Username string `json:"username,omitempty"`
	// Sex 性别
	Sex Sex `json:"sex,omitempty" gorm:"default:0"`
	// IdCode 身份证号码或工号
	IdCode string `json:"idCode,omitempty" gorm:"unique"`
	// Phone 手机号
	Phone string `json:"phone,omitempty" gorm:"unique;not null"`
	// Email 邮箱
	Email string `json:"email,omitempty"`
	// UserDepartments 部门信息
	UserDepartments []*UserDeptInfo `json:"departmentList,omitempty" gorm:"foreignKey:id;joinForeignKey:user_id"`
	// DepartmentsMains 当前用户的部门主管信息
	DepartmentsMains []*UserDeptMain `json:"departmentsMains,omitempty" gorm:"foreignKey:id;joinForeignKey:user_id"`
	// Birthday 生日
	Birthday string `json:"birthday,omitempty"`
	// 员工类型ID
	StaffId string `json:"staffId,omitempty"`
	// Staff 员工类型
	Staff *Staff `json:"staff,omitempty" gorm:"foreignKey:staffId"`
	// Level 员工级别
	Level string `json:"level,omitempty"`
	// JoinTime 加入时间
	JoinTime string `json:"joinTime,omitempty"`
	// TryTime 试用期结束时间
	TryTime string `json:"tryTime,omitempty"`
	// FullTime 转正日期
	FullTime string `json:"fullTime,omitempty"`
	// FirstWorkTime 首次参加工作日期
	FirstWorkTime string `json:"firstWorkTime,omitempty"`
	// NationId 名族ID
	NationId string `json:"nationId,omitempty"`
	// Nation 名族信息
	Nation *Nation `json:"nation,omitempty" gorm:"foreignKey:nationId"`
	// WorkerStatus 在职状态
	WorkerStatus WorkerStatus `json:"workerStatus,omitempty"`
	// Status 状态
	Status UserStatus `json:"status,omitempty"`
	// Comments 备注
	Comments string `json:"comments,omitempty" gorm:"type:longtext"`
	// UpdateAt 数据最后更新日期
	UpdateAt string `json:"updateAt,omitempty"`
	// IsInitManger 是否为初始化管理员
	IsInitManger bool `json:"isInitManger,omitempty"`
	// Icon 头像
	Icon string `json:"icon,omitempty" gorm:"type:longtext;"`
	// CachePassword 用户密码
	CachePassword string `json:"-"`
	// Age 年龄
	Age int `json:"age,omitempty" gorm:"-"`
	// Token token字符串
	Token string `json:"token,omitempty"`
	// RefreshToken 刷新Token
	RefreshToken string `json:"-"`
	// TokenExpire token有效期
	TokenExpire int64 `json:"-"`
}
