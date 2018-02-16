package mathdb

import (
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/mangeshhendre/grpcutils"
	pb "github.com/mangeshhendre/models/services_math_v1"
	_ "github.com/mattn/go-oci8"
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/context"
)

var globalADB *Client
var logger log.Logger

func TestMain(m *testing.M) {
	//objMapperDsn := "acc_crd_conveyor/acc_crd_conveyor@cardsdb00.sgpd.us:1521/crdd01.sgpd.us"
	// your func
	logger = log.New("TestMain")
	dsn := grpcutils.EnvOrDefault("DSN", "userid/password@db_fqdn:1521/service_name?prefetch_rows=5000") //use the connectionstring for database

	DB, err := sqlx.Connect("oci8", dsn)
	if err != nil {
		logger.Fatal("Unable to establish mattn database connection: ", "Error", err)
	}
	DB = DB.Unsafe()

	DB.Mapper = reflectx.NewMapperTagFunc("json", func(str string) string {
		// strToReturn := strings.ToUpper(str)
		// fmt.Println("String", str, "Returning", strToReturn)
		return strings.ToUpper(str)
	},
		func(value string) string {
			//fmt.Println("Value", value)
			var valueToReturn string
			if strings.Contains(value, ",") {
				valueToReturn = strings.Split(value, ",")[0]
			}
			valueToReturn = strings.Replace(valueToReturn, "_", "", -1)
			valueToReturn = strings.ToUpper(valueToReturn)
			//fmt.Println("Value", value, "Returning", valueToReturn)
			return valueToReturn
		})

	// Bypass the cache.
	globalDB, err = New(DB)

	retCode := m.Run()

	defer DB.Close()

	// your func
	//teardown()

	// call with result of m.Run()
	os.Exit(retCode)
}

var mathCases = []struct {
	Case         string
	PropertyID   int64
	OrderNumber  int64
	WantErr      bool
	ObservedDate string
}{
	{
		Case:    "No Results",
		Number1: 0,
		Number2: 0,
		WantErr: true,
	},
	{
		Case:    "Good",
		Number1: 15266709,
		Number2: 600015141,
		WantErr: false,
	},
}

//func TestClient_GetWorkOrderDate(t *testing.T) {
//	foo, err := globalADB.getWorkOrderDate((600015141))
//	if err != nil {
//		t.Error("Unable to get work order date", err)
//	}
//	t.Logf("Sec: %d, nsec: %d , loc: %#v", foo.Second(), foo.Nanosecond(), foo.Location())
//}

func TestClient_AddNumber(t *testing.T) {
	for n, c := range mathCases {
		if c.Number1 == 0 || c.Number2 == 0 {
			continue
		}
		request := &pb.MathRequest{
			Number1: c.PropertyID,
			Number2: c.OrderNumber,
		}
		_, err := globalDB.AddNumber(context.TODO(), request)
		if err != nil {
			if !c.WantErr {
				t.Errorf("Case: %d: %s: Unable to get work order date. Unexpected Error: %s", n, c.Case, err.Error())
			}
			continue
		} else {
			if c.WantErr {
				t.Errorf("Case: %d: %s: Unable to get work order date. Expected error, none reported.", n, c.Case)
				continue
			}
		}
	}
}

func TestClient_MultiplyNumber(t *testing.T) {
	for n, c := range mathCases {
		if c.Number1 == 0 || c.Number2 == 0 {
			continue
		}
		request := &pb.MathRequest{
			Number1: c.PropertyID,
			Number2: c.OrderNumber,
		}
		_, err := globalDB.AddNumber(context.TODO(), request)
		if err != nil {
			if !c.WantErr {
				t.Errorf("Case: %d: %s: Unable to get work order date. Unexpected Error: %s", n, c.Case, err.Error())
			}
			continue
		} else {
			if c.WantErr {
				t.Errorf("Case: %d: %s: Unable to get work order date. Expected error, none reported.", n, c.Case)
				continue
			}
		}
	}
}

func TestClient_DevideNumber(t *testing.T) {
	for n, c := range mathCases {
		if c.Number1 == 0 || c.Number2 == 0 {
			continue
		}
		request := &pb.MathRequest{
			Number1: c.PropertyID,
			Number2: c.OrderNumber,
		}
		_, err := globalDB.AddNumber(context.TODO(), request)
		if err != nil {
			if !c.WantErr {
				t.Errorf("Case: %d: %s: Unable to get work order date. Unexpected Error: %s", n, c.Case, err.Error())
			}
			continue
		} else {
			if c.WantErr {
				t.Errorf("Case: %d: %s: Unable to get work order date. Expected error, none reported.", n, c.Case)
				continue
			}
		}
	}
}
