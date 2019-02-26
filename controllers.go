package controllers
 
import (
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    "kocard/services"
    "kocard/models"
    "time"
)
 
// InvoiceController operations for Invoice
type InvoiceController struct {
    beego.Controller
}
 
type adminTradeView struct {
    Id             uint
    UserCardId     uint
    Price          uint
    Alias          uint
    CurrencyId     uint
    Status         string
    CollectionAlias          string
    Quantity       int
    ImgPath        string
    ImgReversePath string
}
 
type floorTrView struct {
    CardId         uint
    CollectionId   uint
    CardHash       string
    CardStatus     string
    Comment        string
    Alias          string
    UserId         uint
    InvoiceId      uint
    InvoiceStatus  string
    CollectionAlias string
    Price          uint
    Fee          uint
    CurrencyId     uint
    CurrSign       string
}
 
type floorUserTrView struct {
    Id          uint
    CreatedAt   time.Time
    Price           uint
    Currency        string
    CurrencyId      uint
    Status          string
    CollectionAlias  string
    ImgPath        string
    ImgReversePath string
}
 
type InvListView struct {
    CardId         uint
    CollectionId   uint
    CardHash       string
    CardStatus     string
    Comment        string
    Quantity       string
    CllAlias       string
    ImgPath        string
    ImgReversePath string
    Alias          string
    UserId         uint
    InvoiceId      uint
    InvoiceStatus  string
    Price          uint
    Fee            uint
    CurrencyId     uint
    CurrSign       string
}
 
type UserInvView struct {
    UserId          uint
    UserEmail       string
    UserAlias       string
    CollectionAlias string
    CardHash        string
    CardStatus      string
    Quantity        string
    ImgPath         string
    ImgReversePath  string
    InvoiceId       uint
    InvoiceStatus   string
    Price           uint
    CurrencyId      uint
    CurrencySign        string
    InvoiceCreatedAt string
    InvoiceUpdatedAt string
}
 
func getAdmId() int {
    o := orm.NewOrm()
    o.Using("default")
    var user models.Users
    err := o.QueryTable("users").Filter("role", 0).Limit(1).One(&user)
 
    if err != nil {
        return 0
    }
 
    return user.Id
}
 
func (this *InvoiceController) ListInvoices() {
    o := orm.NewOrm()
    o.Using("default")
 
    admId := getAdmId()
    if admId == 0 {
        this.responseWithError(500, map[string]string{"message": "No admin"}, nil)
 
        return
    }
 
    var InvList []InvListView
 
    _, err := o.Raw("select min(cd.collection_num) card_id, collection_id, min(card_hash), cd.status card_status, " +
        "cd.comment, cll.quantity, cll.alias cll_alias, cll.img_path, cll.img_reverse_path, us.alias, us.id user_id, " +
        "min(inv.id) invoice_id,inv.status invoice_status, inv.price, 0 fee, inv.currency_id, curr.sign curr_sign, " +
        "min(cd.id) from cards cd inner join collections cll ON (cd.collection_id=cll.id) " +
        "inner join user_cards uc on (cd.id=uc.card_id) " +
        "inner join users us on (uc.user_id=us.id) " +
        "inner join invoices inv on (uc.id=inv.user_card_id and inv.status='created') " +
        "inner join currencies curr on (curr.id=inv.currency_id) " +
        "where uc.user_id=? and cd.status='minted' group by cd.collection_id, cd.status, " +
            "cd.comment, cll.quantity, cll.alias, cll.img_path, cll.img_reverse_path, us.alias, us.id, " +
            "inv.status, inv.price, inv.currency_id, curr.sign",admId).QueryRows(&InvList)
 
    if err != nil {
        this.responseWithError(500, map[string]string{"message": err.Error()}, err)
 
        return
    }
 
    this.Data["json"] = map[string]interface{}{
        "result": InvList,
    }
 
    this.ServeJSON()
    this.StopRun()
}
 
type UserFloortrView struct {
    TranId        uint
    TranType      string
    InputWalletId uint
    OuputWalletId uint
    CurrencyId    uint
    CurrSign      string
    Amount        uint
    EnyType       string
    EnyStatus     string
}
 
func (this *InvoiceController) UserFloortr() {
    data := this.Ctx.Input.Data()
    email := services.Trim(data["email"].(string))
 
    o := orm.NewOrm()
    o.Using("default")
 
    var user models.Users
    err := o.QueryTable("users").Filter("email", email).Limit(1).One(&user)
 
    if err != nil {
        this.responseWithError(500, map[string]string{"message": err.Error()}, err)
 
        return
    }
 
    /*
        CardId         uint
        CollectionId   uint
        CardHash       string
        CardStatus     string
        Comment        string
        Alias          string
        UserId         uint
        InvoiceId      uint
        InvoiceStatus  string
        Price          uint
        CurrencyId     uint
        CurrSign       string
    */
    var floorTr []floorUserTrView
 
    _, err = o.Raw("SELECT inv.id, inv.created_at, inv.price, cr.alias currency, inv.currency_id, inv.status, cll.alias collection_alias, cll.quantity, cll.img_path " +
    "from invoices inv " +
    "INNER JOIN user_cards uc ON (inv.user_card_id = uc.id) " +
    "INNER JOIN currencies cr on (cr.id=inv.currency_id) " +
    "INNER JOIN cards cd ON uc.card_id = cd.id " +
    "INNER JOIN collections cll ON cd.collection_id = cll.id " +
    "where inv.status='created' and uc.user_id = ?", user.Id).QueryRows(&floorTr)
 
    this.Data["json"] = map[string]interface{}{
        "result": floorTr,
    }
 
    this.ServeJSON()
    this.StopRun()
}
 
func (this *InvoiceController) Floortr() {
    o := orm.NewOrm()
    o.Using("default")
 
    admId := getAdmId()
    if admId == 0 {
        this.responseWithError(500, map[string]string{"message": "No admin"}, nil)
 
        return
    }
 
    collects := []models.Collections{}
    numColl, err := o.QueryTable("collections").All(&collects)
 
    collectRows := map[int]interface{}{}
    if err==nil && numColl > 0 {
 
        var floorTr []floorTrView
 
        for _, collect := range collects {
            numCard, err := o.Raw("select cd.collection_num card_id, cd.collection_id, cd.card_hash, cd.status card_status, " +
                "cd.comment, us.alias, us.id user_id, inv.id invoice_id,inv.status invoice_status, cll.alias collection_alias, inv.price, 0 fee, " +
                "inv.currency_id, curr.sign curr_sign from cards cd  " +
                "inner join collections cll on (cll.id = cd.collection_id)" +
                "inner join user_cards uc on (cd.id=uc.card_id) " +
                "inner join users us on (uc.user_id=us.id) " +
                "inner join invoices inv on (uc.id=inv.user_card_id and inv.status='created') " +
                "inner join currencies curr on (curr.id=inv.currency_id) " +
                "where cd.collection_id = ? and uc.user_id<>?", collect.Id, admId).QueryRows(&floorTr)
 
            if err==nil && numCard > 0 {
 
                collectRows[collect.Id] = map[string]interface{}{
                        "collect" : collect,
                        "cards_num" : numCard,
                        "cards" : floorTr,
 
                }
 
            }
        }
    }
 
    this.Data["json"] = map[string]interface{}{
        "result": collectRows,
    }
 
    this.ServeJSON()
    this.StopRun()
}
 
//FUNC API ADD INVOICE
func (this *InvoiceController) CreateInvoice() {
    data := this.Ctx.Input.Data()
    email := services.Trim(data["email"].(string))
 
    price, errorPrice := this.GetInt("price")
 
    if errorPrice != nil {
        this.responseWithError(500, map[string]string{"message": errorPrice.Error()}, errorPrice)
 
        return
    }
 
    if usCard, errorUsc := this.GetInt("user_card_id"); errorUsc == nil {
 
        o := orm.NewOrm()
        o.Using("default")
 
        var user models.Users
        err := o.QueryTable("users").Filter("email", email).Limit(1).One(&user)
 
        if err != nil {
            this.responseWithError(500, map[string]string{"message": err.Error()}, err)
 
            return
        }
 
        var userCard models.User_cards
        err = o.QueryTable("user_cards").Filter("user_id", user.Id).Filter("card_id", usCard).Limit(1).One(&userCard)
 
        if err != nil || userCard.Id == 0 || userCard.Active == false {
            this.responseWithError(500, map[string]string{"message": err.Error()}, err)
 
            return
        }
 
        res, err := o.Raw("UPDATE invoices SET status='expired' where status='created' and user_card_id = ?", userCard.Id).Exec()
        if err == nil {
            num, _ := res.RowsAffected()
            beego.Info(num)
        } else {
            beego.Error(err)
        }
 
        inv := new(models.Invoices)
        inv.UserCardId = userCard.Id
        inv.Price = price
        inv.CurrencyId = 1
        inv.Status = "created"
        inv.StatusComment = "Selling Card"
        invId, err := o.Insert(inv)
 
        userCard.InvoiceId = int(invId)
        _,err = o.Update(&userCard, "InvoiceId")
 
        if err == nil {
            this.Data["json"] = map[string]interface{}{
                "result": "Create Invoice",
                "trans_id": invId,
            }
            this.ServeJSON()
            this.StopRun()
        } else {
            this.responseWithError(500, map[string]string{"message": errorUsc.Error()}, errorUsc)
 
            return
        }
 
    } else {
        this.responseWithError(500, map[string]string{"message": errorUsc.Error()}, errorUsc)
 
        return
    }
}
 
//  API error response
func (this *InvoiceController) responseWithError(status int, message map[string]string, err interface{}) {
    beego.Error(err)
 
    this.Ctx.Output.SetStatus(status)
    this.Data["json"] = message
    this.ServeJSON()
    this.StopRun()
 
    return
}
 
// URLMapping ...
func (ic *InvoiceController) URLMapping() {
    ic.Mapping("Post", ic.Post)
    ic.Mapping("GetOne", ic.GetOne)
    ic.Mapping("GetAll", ic.GetAll)
    ic.Mapping("Put", ic.Put)
    ic.Mapping("Delete", ic.Delete)
}
 
// Post ...
// @Title Create
// @Description create Invoice
// @Param   body        body    models.Invoice  true        "body for Invoice content"
// @Success 201 {object} models.Invoice
// @Failure 403 body is empty
// @router / [post]
func (ic *InvoiceController) Post() {
 
}
 
// GetOne ...
// @Title GetOne
// @Description get Invoice by id
// @Param   id      path    string  true        "The key for staticblock"
// @Success 200 {object} models.Invoice
// @Failure 403 :id is empty
// @router /:id [get]
func (ic *InvoiceController) GetOne() {
 
}
 
// GetAll ...
// @Title GetAll
// @Description get Invoice
// @Param   query   query   string  false   "Filter. e.g. col1:v1,col2:v2 ..."
// @Param   fields  query   string  false   "Fields returned. e.g. col1,col2 ..."
// @Param   sortby  query   string  false   "Sorted-by fields. e.g. col1,col2 ..."
// @Param   order   query   string  false   "Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param   limit   query   string  false   "Limit the size of result set. Must be an integer"
// @Param   offset  query   string  false   "Start position of result set. Must be an integer"
// @Success 200 {object} models.Invoice
// @Failure 403
// @router / [get]
func (ic *InvoiceController) GetAll() {
    var i []UserInvView
    o := orm.NewOrm()
    o.Using("default")
 
    _, err := o.Raw("SELECT us.id user_id, us.email user_email, us.alias user_alias, cll.alias collection_alias," +
        "cr.card_hash,cr.status card_status, cll.quantity, cll.img_path," +
        "cll.img_reverse_path, inv.id invoice_id, inv.status invoice_status, inv.price," +
        "curr.id currency_id, curr.sign currency_sign, inv.created_at invoice_created_at, inv.updated_at invoice_updated_at " +
        "from invoices inv " +
        "INNER JOIN currencies curr ON inv.currency_id = curr.id " +
        "INNER JOIN user_cards uc on uc.id = inv.user_card_id " +
        "INNER JOIN cards cr on uc.card_id = cr.id " +
        "INNER JOIN users us on us.id = uc.user_id " +
        "INNER JOIN collections cll ON cr.collection_id = cll.id").QueryRows(&i)
    if err != nil {
        ic.responseWithError(500, map[string]string{"message": err.Error()}, err)
 
        return
    }
 
    ic.Data["json"] = map[string]interface{}{
        "invoices": i,
    }
 
    ic.ServeJSON()
    ic.StopRun()
}
 
// Put ...
// @Title Put
// @Description update the Invoice
// @Param   id      path    string  true        "The id you want to update"
// @Param   body        body    models.Invoice  true        "body for Invoice content"
// @Success 200 {object} models.Invoice
// @Failure 403 :id is not int
// @router /:id [put]
func (ic *InvoiceController) Put() {
 
}
 
// Delete ...
// @Title Delete
// @Description delete the Invoice
// @Param   id      path    string  true        "The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (ic *InvoiceController) Delete() {
 
}