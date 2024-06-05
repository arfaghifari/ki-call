package converter

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyStruct1 struct {
	Field1  string       `json:"field1"`
	Field2  int64        `json:"field2"`
	Field3  *int         `json:"field3"`
	Field4  []string     `json:"field4"`
	Field5  MyStruct2    `json:"field5"`
	Field6  *MyStruct2   `json:"field6"`
	Field7  []MyStruct2  `json:"field7"`
	Field8  []*MyStruct2 `json:"field8"`
	Field9  *int         `json:"field9"`
	Field10 []int        `json:"field10"`
}

type MyStruct2 struct {
	SubFiled1 []*string `json:"sub_field1"`
}

var myWord1 = "satu"
var myWord2 = "dua"
var myWord3 = "tiga"
var myWord4 = "empat"
var myInt1 = 10
var nol = 0
var emptyWord = ""
var subStruct1 = MyStruct2{
	[]*string{&myWord1},
}
var subStruct2 = MyStruct2{
	[]*string{&myWord2},
}
var subStruct3 = MyStruct2{
	[]*string{&myWord3},
}
var subStruct4 = MyStruct2{
	[]*string{&myWord4},
}

var myStruct = MyStruct1{
	Field1:  "1",
	Field2:  1,
	Field3:  &myInt1,
	Field4:  []string{"2", "3"},
	Field5:  subStruct1,
	Field6:  &subStruct2,
	Field7:  []MyStruct2{subStruct3},
	Field8:  []*MyStruct2{&subStruct4},
	Field9:  nil,
	Field10: []int{},
}

var myMap = map[string]interface{}{
	"field1": "1",
	"field2": int64(1),
	"field3": 10,
	"field4": []string{"2", "3"},
	"field5": map[string]interface{}{
		"sub_field1": []string{"satu"},
	},
	"field6": map[string]interface{}{
		"sub_field1": []string{"dua"},
	},
	"field7": []map[string]interface{}{
		{
			"sub_field1": []string{"tiga"},
		},
	},
	"field8": []map[string]interface{}{
		{
			"sub_field1": []string{"empat"},
		},
	},
	"field9":  nil,
	"field10": []int{},
}

var myStructNoEmpty = MyStruct1{
	Field1: "1",
	Field2: 1,
	Field3: &myInt1,
	Field4: []string{"2", "3"},
	Field5: subStruct1,
	Field6: &subStruct2,
	Field7: []MyStruct2{
		{SubFiled1: []*string{&emptyWord}},
	},
	Field8:  []*MyStruct2{&subStruct4},
	Field9:  &nol,
	Field10: []int{0},
}

var myStructNoEmptyDecoy = MyStruct1{
	Field1:  "1",
	Field2:  1,
	Field3:  &myInt1,
	Field4:  []string{"2", "3"},
	Field5:  subStruct1,
	Field6:  &subStruct2,
	Field7:  []MyStruct2{},
	Field8:  []*MyStruct2{&subStruct4},
	Field9:  nil,
	Field10: []int{},
}

var myMapNoEmpty = map[string]interface{}{
	"field1": "1",
	"field2": int64(1),
	"field3": 10,
	"field4": []string{"2", "3"},
	"field5": map[string]interface{}{
		"sub_field1": []string{"satu"},
	},
	"field6": map[string]interface{}{
		"sub_field1": []string{"dua"},
	},
	"field7": []map[string]interface{}{
		{
			"sub_field1": []string{""},
		},
	},
	"field8": []map[string]interface{}{
		{
			"sub_field1": []string{"empat"},
		},
	},
	"field9":  0,
	"field10": []int{0},
}

func TestTransformStructToMapJson(t *testing.T) {
	type args struct {
		val     interface{}
		noEmpty bool
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "no empty false",
			args: args{
				val: myStruct,
			},
			want: myMap,
		},
		{
			name: "no empty true",
			args: args{
				val:     myStructNoEmptyDecoy,
				noEmpty: true,
			},
			want: myMapNoEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformStructToMapJson(tt.args.val, tt.args.noEmpty)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformStructToMapJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			expectedJson, _ := json.Marshal(tt.want)
			myJson, _ := json.Marshal(got)
			if !reflect.DeepEqual(string(myJson), string(expectedJson)) {
				t.Errorf("TransformStructToMapJson() = %v, want %v", string(myJson), string(expectedJson))
			}
		})
	}
}

func TestTransformMapJsonToStruct(t *testing.T) {
	type args struct {
		myMapInterface interface{}
		valueType      reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "case 1",
			args: args{
				myMapInterface: myMap,
				valueType:      reflect.TypeOf(myStruct),
			},
			want: myStruct,
		},
		{
			name: "case 2",
			args: args{
				myMapInterface: myMapNoEmpty,
				valueType:      reflect.TypeOf(myStruct),
			},
			want: myStructNoEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapArgs := map[string]interface{}{}
			expectedJson, _ := json.Marshal(tt.args.myMapInterface)

			_ = json.Unmarshal(expectedJson, &mapArgs)

			got, err := TransformMapJsonToStruct(mapArgs, tt.args.valueType)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformMapJsonToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformMapJsonToStruct() = %v, want %v", got, tt.want)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
