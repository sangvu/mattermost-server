package searchengine

import (
	"github.com/mattermost/mattermost-server/jobs"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/services/searchengine/bleveengine"
	"github.com/mattermost/mattermost-server/services/searchengine/nullengine"
)

func NewSearchEngineBroker(cfg *model.Config, jobServer *jobs.JobServer) (*SearchEngineBroker, error) {
	broker := &SearchEngineBroker{
		cfg:       cfg,
		jobServer: jobServer,
	}

	if *cfg.BleveSettings.EnableIndexing && *cfg.BleveSettings.IndexDir != "" {
		bleveEngine, err := bleveengine.NewBleveEngine(cfg, jobServer)
		if err != nil {
			return nil, err
		}
		broker.BleveEngine = bleveEngine
	}

	nullEngine, err := nullengine.NewNullEngine()
	if err != nil {
		return nil, err
	}
	broker.NullEngine = nullEngine
	return broker, nil
}

func (seb *SearchEngineBroker) RegisterElasticsearchEngine(es SearchEngineInterface) {
	seb.ElasticsearchEngine = es
}

type SearchEngineBroker struct {
	cfg                 *model.Config
	jobServer           *jobs.JobServer
	BleveEngine         SearchEngineInterface
	ElasticsearchEngine SearchEngineInterface
	NullEngine          SearchEngineInterface
}

func (seb *SearchEngineBroker) UpdateConfig(cfg *model.Config) *model.AppError {
	seb.cfg = cfg
	return nil
}

func (seb *SearchEngineBroker) GetActiveEngines() []SearchEngineInterface {
	engines := []SearchEngineInterface{}
	if seb.ElasticsearchEngine != nil && seb.ElasticsearchEngine.IsActive() {
		engines = append(engines, seb.ElasticsearchEngine)
	}
	if seb.BleveEngine != nil && seb.BleveEngine.IsActive() {
		engines = append(engines, seb.BleveEngine)
	}
	engines = append(engines, seb.NullEngine)
	return engines
}