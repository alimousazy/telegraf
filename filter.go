package telegraf


type Filter interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Description returns a one-sentence description on the Input
	Description() string

  //Output metric to outputs list
  OutputMetric(output interface{})
  //Add metric to the middleware
  AddMetric(metric Metric)
  //Called on each metric to check if this middle ware enabled
  //or not for that metric 
  IsEnabled(name string) bool
  //clear metrics to output 
  Reset()
}
