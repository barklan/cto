package storage

var floatType = "float64"

type VariableKeys struct {
	ErrSimThresh []string
}

var KVars = VariableKeys{
	ErrSimThresh: []string{"err.sim.thresh", floatType},
}

func (data *Data) KVarExists(key string) bool {
	return false
}

func (data *Data) SetKvarValue() {}
