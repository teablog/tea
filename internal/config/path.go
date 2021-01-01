package config

var Path *_path

type _path struct{}

func (*_path) StorageDir() string {
	return Config.Section("path").Key("storage_dir").String()
}
