package service

type SettingDTO struct {
	GachaCoinConsumption int `json:"gachaCoinConsumption"`
}

type SettingService interface {
	GetSetting() (*SettingDTO, error)
}

type settingService struct {
	// ここに必要なリポジトリを追加
}

func NewSettingService() SettingService {
	return &settingService{
		// ここで必要なリポジトリを初期化
	}
}

const (
	// ガチャ1回あたりのコイン消費量
	GachaCoinConsumption = 100
)

func (s *settingService) GetSetting() (*SettingDTO, error) {
	// ここでは固定値を返すが、将来的にはDBや設定ファイルから取得するよう変更可能にする
	return &SettingDTO{
		GachaCoinConsumption: GachaCoinConsumption,
	}, nil
}
