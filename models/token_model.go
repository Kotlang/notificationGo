package models

type DeviceInstanceModel struct {
	LoginId string `bson:"_id" json:"loginId"`
	Token   string `bson:"token" json:"token"`
	Tenant  string `bson:"tenant" json:"tenant"`
}

func (m *DeviceInstanceModel) Id() string {
	return m.LoginId
}
