package logmon 

import (
	ejson "encoding/json"
  "fmt"
	"github.com/influxdata/telegraf"
  "regexp"
  "strings"
  "strconv"
)


type LogmonSerializer struct {
}


func (s *LogmonSerializer) Serialize(metric telegraf.Metric) ([]string, error) {
  pctRegex := regexp.MustCompile("p([0-9]+)")
	out := []string{}
	fields := []string{}
  tags  := metric.Tags()
	m := make(map[string]map[string]string)

  for k, v := range metric.Fields() {
    value, ok := v.(float64); if !ok {
      continue
    }
    stValue := strconv.FormatFloat(value, 'f', -1, 64)
    pctRegex.ReplaceAllString(k, "pct$1")
    list := strings.Split(k, ".")
    maping := map[string]string{
      "count" : "cnt",
    }
    if(len(list) == 2) {
      mName := metric.Name() + "." + list[0]
      if m[mName] == nil {
        m[mName] = make(map[string]string)
        m[mName]["type"] = "HIS"
      }
      t, ok := maping[list[1]] ; if (!ok) {
        m[mName][list[1]] = stValue
      } else  {
        m[mName][t] = stValue
      }
    } else {
      mName := metric.Name() + "." + k
      _, ok := m[mName]; if (!ok) {
         m[mName] = make(map[string]string)
      }
      m[mName]["value"] = stValue
    }
  }
	serialized, err := ejson.Marshal(m)
	if err != nil {
		return []string{}, err
	}
  year, month, day := metric.Time().Date()
  hour, min, sec := metric.Time().Clock()
  
//WMPLTFMLOG254103  1466690220371 2016-06-23 13:57:00.371 rules-api-11030078-16-54937210  - 54937210  - - PROD  ship-pricing-rules  prod-dfw3 prod  3.1.27  9f312313-1f-1557d8cb153005  ME  not_applicable  - pool_ship_pricing_gecwalmart_comussprp1ship_pricing_gecwalmart_com  activate  - {"absolute":{"unit":"Calls","min":0,"soq":0,"max":0,"cnt":0,"sum":0,"type":"BAS"},"TIMER":{"unit":"us","min":0,"soq":0,"pct95":0,"max":0,"pct75":0,"cnt":0,"sum":0,"type":"HIS","pct999":0}}
  fields = append(fields, tags["wm_sign"])
  fields = append(fields, strconv.FormatInt(metric.UnixNano() / 1000, 10))
  fields = append(fields, fmt.Sprintf("%02d-%02d-%02d", year, int(month), day))
  fields = append(fields, fmt.Sprintf("%02d:%02d:%02d.000", hour, min, sec))
  fields = append(fields, tags["host"])
  fields = append(fields, "-")
  fields = append(fields, tags["compute_id"])
  fields = append(fields, "-")
  fields = append(fields, "-")
  fields = append(fields, tags["envtype"])
  fields = append(fields, tags["app_name"])
  fields = append(fields, tags["dc"])
  fields = append(fields, tags["env"])
  fields = append(fields, tags["app_version"])
  fields = append(fields, tags["msg_id"])
  fields = append(fields, "ME")
  fields = append(fields, tags["level1"])
  fields = append(fields, "-")
  fields = append(fields, tags["level2"])
  fields = append(fields, tags["level3"])
  fields = append(fields, "-")
	fields = append(fields, string(serialized))

  out = append(out, strings.Join(fields, "\t"))
	return out, nil
}
