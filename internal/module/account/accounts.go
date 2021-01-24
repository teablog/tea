package account

type AcctSlice []*Account

func (rows AcctSlice) NameRepeat(name, email string) bool {
	for _, v := range rows {
		if v.Name == name && v.Email != email {
			return true
		}
	}
	return false
}

func (rows AcctSlice) Acct(name, email string) (*Account, bool) {
	for _, v := range rows {
		if v.Name == name && v.Email == email {
			return v, true
		}
	}
	return nil, false
}
