package service

type SettingDTO struct {
	GachaCoinConsumption int `json:"gachaCoinConsumption"`
	RankingFetchCount    int `json:"rankingFetchCount"`
}

type SettingService interface {
	Get() (*SettingDTO, error)
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
	// ランキングを取得する際の取得件数
	RankingFetchCount = 10
)

func (s *settingService) Get() (*SettingDTO, error) {
	// ここでは固定値を返すが、将来的にはDBや設定ファイルから取得するよう変更可能にする
	return &SettingDTO{
		GachaCoinConsumption: GachaCoinConsumption,
		RankingFetchCount:    RankingFetchCount,
	}, nil
}
