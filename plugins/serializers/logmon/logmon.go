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


func (s *LogmonSerializer) inArray(list []string, item string) (bool) {
  for _, l := range list {
    if(strings.Compare(l, item) == 0) {
      return true;
    }
  }
  return false;
}
func (s *LogmonSerializer) Serialize(metric telegraf.Metric) ([]string, error) {
  pctRegex := regexp.MustCompile("p([0-9]+)")
	out := []string{}
	fields := []string{}
  tags  := metric.Tags()
  ignore_field := []string{"variance"}
	m := make(map[string]map[string]interface{})

  for k, v := range metric.Fields() {
    value, ok := v.(float64); if !ok {
      continue
    }
    k = pctRegex.ReplaceAllString(k, "pct$1")
    list := strings.Split(k, ".")
    maping := map[string]string{
      "count" : "cnt",
      "mean" : "soq",
    }
    if(len(list) == 2) {
      mName := metric.Name() + "." + list[0]
      if m[mName] == nil {
        m[mName] = make(map[string]interface{})
        m[mName]["type"] = "HIS"
        m[mName]["unit"] = "n"
      }
      if s.inArray(ignore_field, list[1]) {
        continue
      }
      t, ok := maping[list[1]] ; if (!ok) {
        m[mName][list[1]] = value
      } else  {
        m[mName][t] = value
      }
    } else {
      //"cur":0,"unit":"Calls","min":0,"max":0,"type":"GAU"
      mName := metric.Name() + "." + k
      _, ok := m[mName]; if (!ok) {
         m[mName] = make(map[string]interface{})
      }
      m[mName]["cur"] = value
      m[mName]["min"] = value
      m[mName]["max"] = value
      m[mName]["unit"] = "n"
      m[mName]["type"] = "GAU"
    }
  }
	serialized, err := ejson.Marshal(m)
	if err != nil {
		return []string{}, err
	}
  year, month, day := metric.Time().Date()
  hour, min, sec := metric.Time().Clock()
  fields = append(fields, tags["wm_sign"])
  fields = append(fields, strconv.FormatInt(metric.UnixNano() / 1000000, 10))
  fields = append(fields, fmt.Sprintf("%02d-%02d-%02d", year, int(month), day) + " " + fmt.Sprintf("%02d:%02d:%02d.000", hour, min, sec))
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
