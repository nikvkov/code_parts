package address
 
import (
    "encoding/json"
    _ "fmt"
    "github.com/jinzhu/gorm"
    "github.com/jmoiron/sqlx"
    "github.com/labstack/echo"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "io/ioutil"
    _ "log"
    "net/http"
    "net/http/httptest"
    "testing"
    "workload/jiazhen-api/pkg/base"
    "workload/jiazhen-api/pkg/database"
    "workload/jiazhen-api/pkg/jwt"
    "strings"
)
 
type TestSuite struct {
    suite.Suite
 
    dbName string
    gormDB *gorm.DB
    sqlDB  *sqlx.DB
}
 
func (s *TestSuite) SetupSuite() {
    var err error
 
    s.dbName, err = database.CreateDB()
    if err != nil {
        s.T().Errorf("error creating test database: %v", err)
    }
 
    s.gormDB, err = database.NewGorm(s.dbName)
    if err != nil {
        s.T().Errorf("error connecting to gorm db: %v", err)
    }
}
 
func (s *TestSuite) TearDownSuite() {
    err := database.DropDB(s.dbName)
    if err != nil {
        s.T().Errorf("could not drop database %v: %v", s.dbName, err)
    }
}
 
func TestTestSuiteIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipped integration test")
    }
    suite.Run(t, new(TestSuite))
}
 
type TestUserWechatUserAddresses struct {
    Status int             `json:"status"`
    Data   []ParseAddr `json:"data"`
}
 
func (s *TestSuite) TestGetUserAddressList() {
    // Test wechat get user addresses list.
    m := NewManager(s.sqlDB, s.gormDB)
    e := echo.New()
 
    // Test the gettig user addresses
    req := httptest.NewRequest(echo.GET, "/wechat/user/addresses", nil)
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.Set("user", jwt.GenTestUserJwt(1))
 
    if assert.NoError(s.T(), m.GetUserAddressList(c)) {
        assert.Equal(s.T(), http.StatusOK, rec.Code)
 
        // Read the body.
        body, err := ioutil.ReadAll(rec.Body)
        assert.Equal(s.T(), nil, err)
 
        // Convert the body to Go struct.
        r := &TestUserWechatUserAddresses{}
        err = json.Unmarshal(body, r)
        assert.Equal(s.T(), nil, err)
 
        assert.Equal(s.T(), base.GeneralStatusSuccess, r.Status)
        data := r.Data
        assert.Equal(s.T(), 1, len(data))
        assert.Equal(s.T(), 1, data[0].ID)
        assert.Equal(s.T(), "123456789098", data[0].Phone)
        assert.Equal(s.T(), "userName1", data[0].UserName)
        assert.Equal(s.T(), "userName1 detail address", data[0].Detail)
    }
}
 
func (s *TestSuite) TestUpdateUserAddress () {
    // Test wechat update user address.
    m := NewManager(s.sqlDB, s.gormDB)
    e := echo.New()
 
    // Test the updating user address
    reqBody := strings.NewReader(`{"UserName": "UserNameUpdate", "Phone": "333333333333", "Detail": "UserName address detail update"}`)
    reqUp := httptest.NewRequest("POST", "/wechat/user/address/1", reqBody)
    reqUp.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    c := e.NewContext(reqUp, w)
 
    // This is essential for setting the path.
    c.SetPath("/wechat/user/address/:address_id")
    c.SetParamNames("address_id")
    c.SetParamValues("1")
 
    c.Set("user", jwt.GenTestUserJwt(1))
 
    if assert.NoError(s.T(), m.UpdateUserAddress(c)) {
        assert.Equal(s.T(), http.StatusOK, w.Code)
 
        req := httptest.NewRequest(echo.GET, "/wechat/user/addresses", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", jwt.GenTestUserJwt(1))
 
        if assert.NoError(s.T(), m.GetUserAddressList(c)) {
 
            assert.Equal(s.T(), http.StatusOK, rec.Code)
            // Read the body.
            body, err := ioutil.ReadAll(rec.Body)
            assert.Equal(s.T(), nil, err)
 
            // Convert the body to Go struct.
            r := &TestUserWechatUserAddresses{}
            err = json.Unmarshal(body, r)
            assert.Equal(s.T(), nil, err)
 
            assert.Equal(s.T(), base.GeneralStatusSuccess, r.Status)
            data := r.Data
 
            assert.Equal(s.T(), 1, len(data))
            assert.Equal(s.T(), 1, data[0].ID)
            assert.Equal(s.T(), "333333333333", data[0].Phone)
            assert.Equal(s.T(), "UserNameUpdate", data[0].UserName)
            assert.Equal(s.T(), "UserName address detail update", data[0].Detail)
        }
    }
}
 
func (s *TestSuite) TestCreateUserAddress () {
    // Test wechat create user address.
    m := NewManager(s.sqlDB, s.gormDB)
    e := echo.New()
 
    // Test the creating user address
    reqBody := strings.NewReader(`{"UserName": "UserNameInsert", "Phone": "222222222222", "Detail": "UserName address detail"}`)
    reqIns := httptest.NewRequest("POST", "/wechat/user/address", reqBody)
    reqIns.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    c := e.NewContext(reqIns, w)
    c.Set("user", jwt.GenTestUserJwt(4))
 
    if assert.NoError(s.T(), m.CreateUserAddress(c)) {
        assert.Equal(s.T(), http.StatusOK, w.Code)
 
        req := httptest.NewRequest(echo.GET, "/wechat/user/addresses", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", jwt.GenTestUserJwt(4))
 
        if assert.NoError(s.T(), m.GetUserAddressList(c)) {
            assert.Equal(s.T(), http.StatusOK, rec.Code)
            // Read the body.
            body, err := ioutil.ReadAll(rec.Body)
            assert.Equal(s.T(), nil, err)
 
            // Convert the body to Go struct.
            r := &TestUserWechatUserAddresses{}
            err = json.Unmarshal(body, r)
            assert.Equal(s.T(), nil, err)
 
            assert.Equal(s.T(), base.GeneralStatusSuccess, r.Status)
            data := r.Data
 
            assert.Equal(s.T(), 1, len(data))
            assert.Equal(s.T(), 4, data[0].ID)
            assert.Equal(s.T(), 4, data[0].UserID)
            assert.Equal(s.T(), "222222222222", data[0].Phone)
            assert.Equal(s.T(), "UserNameInsert", data[0].UserName)
            assert.Equal(s.T(), "UserName address detail", data[0].Detail)
        }
    }
}
 
func (s *TestSuite) TestDeleteUserAddress () {
    // Test wechat delete user address.
    m := NewManager(s.sqlDB, s.gormDB)
    e := echo.New()
 
    // Test the deleting user address
    reqD := httptest.NewRequest("DELETE", "/wechat/user/address/4", nil)
    reqD.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    c := e.NewContext(reqD, w)
 
    // This is essential for setting the path.
    c.SetPath("/wechat/user/address/:address_id")
    c.SetParamNames("address_id")
    c.SetParamValues("4")
 
    c.Set("user", jwt.GenTestUserJwt(4))
 
    if assert.NoError(s.T(), m.DeleteUserAddress(c)) {
        assert.Equal(s.T(), http.StatusOK, w.Code)
 
        req := httptest.NewRequest(echo.GET, "/wechat/user/addresses", nil)
        req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.Set("user", jwt.GenTestUserJwt(4))
 
        if assert.NoError(s.T(), m.GetUserAddressList(c)) {
 
            assert.Equal(s.T(), http.StatusOK, rec.Code)
            // Read the body.
            body, err := ioutil.ReadAll(rec.Body)
            assert.Equal(s.T(), nil, err)
 
            // Convert the body to Go struct.
            r := &TestUserWechatUserAddresses{}
            err = json.Unmarshal(body, r)
            assert.Equal(s.T(), nil, err)
 
            assert.Equal(s.T(), base.GeneralStatusSuccess, r.Status)
            data := r.Data
 
            assert.Equal(s.T(), 0, len(data))
        }
    }
}