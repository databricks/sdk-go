package v1

type ForecastingExperiment_State string

const (
	ForecastingExperiment_State_Pending   ForecastingExperiment_State = "PENDING"
	ForecastingExperiment_State_Running   ForecastingExperiment_State = "RUNNING"
	ForecastingExperiment_State_Succeeded ForecastingExperiment_State = "SUCCEEDED"
	ForecastingExperiment_State_Failed    ForecastingExperiment_State = "FAILED"
	ForecastingExperiment_State_Cancelled ForecastingExperiment_State = "CANCELLED"
)

// String representation for [fmt.Print].
func (f *ForecastingExperiment_State) String() string {
	return string(*f)
}

type CreateForecastingExperimentRequest struct {
	TrainDataPath               *string  `json:"train_data_path"`
	TargetColumn                *string  `json:"target_column"`
	TimeColumn                  *string  `json:"time_column"`
	ForecastGranularity         *string  `json:"forecast_granularity"`
	ForecastHorizon             *int64   `json:"forecast_horizon"`
	PrimaryMetric               *string  `json:"primary_metric"`
	TrainingFrameworks          []string `json:"training_frameworks"`
	ExperimentPath              *string  `json:"experiment_path"`
	MaxRuntime                  *int64   `json:"max_runtime"`
	SplitColumn                 *string  `json:"split_column"`
	CustomWeightsColumn         *string  `json:"custom_weights_column"`
	RegisterTo                  *string  `json:"register_to"`
	HolidayRegions              []string `json:"holiday_regions"`
	TimeseriesIdentifierColumns []string `json:"timeseries_identifier_columns"`
	PredictionDataPath          *string  `json:"prediction_data_path"`
	IncludeFeatures             []string `json:"include_features"`
	FutureFeatureDataPath       *string  `json:"future_feature_data_path"`
}

type CreateForecastingExperimentResponse struct {
	ExperimentId *string `json:"experiment_id"`
}

// Represents a forecasting experiment with its unique identifier, URL, and
// state..
type ForecastingExperiment struct {
	ExperimentId      *string                      `json:"experiment_id"`
	ExperimentPageUrl *string                      `json:"experiment_page_url"`
	State             *ForecastingExperiment_State `json:"state"`
}

type GetForecastingExperimentRequest struct {
	ExperimentId *string `json:"experiment_id"`
}
