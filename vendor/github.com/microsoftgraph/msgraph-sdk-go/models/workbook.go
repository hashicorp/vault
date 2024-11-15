package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

type Workbook struct {
    Entity
}
// NewWorkbook instantiates a new Workbook and sets the default values.
func NewWorkbook()(*Workbook) {
    m := &Workbook{
        Entity: *NewEntity(),
    }
    return m
}
// CreateWorkbookFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateWorkbookFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWorkbook(), nil
}
// GetApplication gets the application property value. The application property
// returns a WorkbookApplicationable when successful
func (m *Workbook) GetApplication()(WorkbookApplicationable) {
    val, err := m.GetBackingStore().Get("application")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookApplicationable)
    }
    return nil
}
// GetComments gets the comments property value. Represents a collection of comments in a workbook.
// returns a []WorkbookCommentable when successful
func (m *Workbook) GetComments()([]WorkbookCommentable) {
    val, err := m.GetBackingStore().Get("comments")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookCommentable)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Workbook) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["application"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookApplicationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetApplication(val.(WorkbookApplicationable))
        }
        return nil
    }
    res["comments"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookCommentFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookCommentable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookCommentable)
                }
            }
            m.SetComments(res)
        }
        return nil
    }
    res["functions"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetObjectValue(CreateWorkbookFunctionsFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            m.SetFunctions(val.(WorkbookFunctionsable))
        }
        return nil
    }
    res["names"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookNamedItemFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookNamedItemable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookNamedItemable)
                }
            }
            m.SetNames(res)
        }
        return nil
    }
    res["operations"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookOperationFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookOperationable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookOperationable)
                }
            }
            m.SetOperations(res)
        }
        return nil
    }
    res["tables"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookTableFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookTableable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookTableable)
                }
            }
            m.SetTables(res)
        }
        return nil
    }
    res["worksheets"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWorkbookWorksheetFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WorkbookWorksheetable, len(val))
            for i, v := range val {
                if v != nil {
                    res[i] = v.(WorkbookWorksheetable)
                }
            }
            m.SetWorksheets(res)
        }
        return nil
    }
    return res
}
// GetFunctions gets the functions property value. The functions property
// returns a WorkbookFunctionsable when successful
func (m *Workbook) GetFunctions()(WorkbookFunctionsable) {
    val, err := m.GetBackingStore().Get("functions")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(WorkbookFunctionsable)
    }
    return nil
}
// GetNames gets the names property value. Represents a collection of workbooks scoped named items (named ranges and constants). Read-only.
// returns a []WorkbookNamedItemable when successful
func (m *Workbook) GetNames()([]WorkbookNamedItemable) {
    val, err := m.GetBackingStore().Get("names")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookNamedItemable)
    }
    return nil
}
// GetOperations gets the operations property value. The status of workbook operations. Getting an operation collection is not supported, but you can get the status of a long-running operation if the Location header is returned in the response. Read-only.
// returns a []WorkbookOperationable when successful
func (m *Workbook) GetOperations()([]WorkbookOperationable) {
    val, err := m.GetBackingStore().Get("operations")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookOperationable)
    }
    return nil
}
// GetTables gets the tables property value. Represents a collection of tables associated with the workbook. Read-only.
// returns a []WorkbookTableable when successful
func (m *Workbook) GetTables()([]WorkbookTableable) {
    val, err := m.GetBackingStore().Get("tables")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookTableable)
    }
    return nil
}
// GetWorksheets gets the worksheets property value. Represents a collection of worksheets associated with the workbook. Read-only.
// returns a []WorkbookWorksheetable when successful
func (m *Workbook) GetWorksheets()([]WorkbookWorksheetable) {
    val, err := m.GetBackingStore().Get("worksheets")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.([]WorkbookWorksheetable)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Workbook) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    {
        err = writer.WriteObjectValue("application", m.GetApplication())
        if err != nil {
            return err
        }
    }
    if m.GetComments() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetComments()))
        for i, v := range m.GetComments() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("comments", cast)
        if err != nil {
            return err
        }
    }
    {
        err = writer.WriteObjectValue("functions", m.GetFunctions())
        if err != nil {
            return err
        }
    }
    if m.GetNames() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetNames()))
        for i, v := range m.GetNames() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("names", cast)
        if err != nil {
            return err
        }
    }
    if m.GetOperations() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetOperations()))
        for i, v := range m.GetOperations() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("operations", cast)
        if err != nil {
            return err
        }
    }
    if m.GetTables() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetTables()))
        for i, v := range m.GetTables() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("tables", cast)
        if err != nil {
            return err
        }
    }
    if m.GetWorksheets() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetWorksheets()))
        for i, v := range m.GetWorksheets() {
            if v != nil {
                cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
            }
        }
        err = writer.WriteCollectionOfObjectValues("worksheets", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetApplication sets the application property value. The application property
func (m *Workbook) SetApplication(value WorkbookApplicationable)() {
    err := m.GetBackingStore().Set("application", value)
    if err != nil {
        panic(err)
    }
}
// SetComments sets the comments property value. Represents a collection of comments in a workbook.
func (m *Workbook) SetComments(value []WorkbookCommentable)() {
    err := m.GetBackingStore().Set("comments", value)
    if err != nil {
        panic(err)
    }
}
// SetFunctions sets the functions property value. The functions property
func (m *Workbook) SetFunctions(value WorkbookFunctionsable)() {
    err := m.GetBackingStore().Set("functions", value)
    if err != nil {
        panic(err)
    }
}
// SetNames sets the names property value. Represents a collection of workbooks scoped named items (named ranges and constants). Read-only.
func (m *Workbook) SetNames(value []WorkbookNamedItemable)() {
    err := m.GetBackingStore().Set("names", value)
    if err != nil {
        panic(err)
    }
}
// SetOperations sets the operations property value. The status of workbook operations. Getting an operation collection is not supported, but you can get the status of a long-running operation if the Location header is returned in the response. Read-only.
func (m *Workbook) SetOperations(value []WorkbookOperationable)() {
    err := m.GetBackingStore().Set("operations", value)
    if err != nil {
        panic(err)
    }
}
// SetTables sets the tables property value. Represents a collection of tables associated with the workbook. Read-only.
func (m *Workbook) SetTables(value []WorkbookTableable)() {
    err := m.GetBackingStore().Set("tables", value)
    if err != nil {
        panic(err)
    }
}
// SetWorksheets sets the worksheets property value. Represents a collection of worksheets associated with the workbook. Read-only.
func (m *Workbook) SetWorksheets(value []WorkbookWorksheetable)() {
    err := m.GetBackingStore().Set("worksheets", value)
    if err != nil {
        panic(err)
    }
}
type Workbookable interface {
    Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetApplication()(WorkbookApplicationable)
    GetComments()([]WorkbookCommentable)
    GetFunctions()(WorkbookFunctionsable)
    GetNames()([]WorkbookNamedItemable)
    GetOperations()([]WorkbookOperationable)
    GetTables()([]WorkbookTableable)
    GetWorksheets()([]WorkbookWorksheetable)
    SetApplication(value WorkbookApplicationable)()
    SetComments(value []WorkbookCommentable)()
    SetFunctions(value WorkbookFunctionsable)()
    SetNames(value []WorkbookNamedItemable)()
    SetOperations(value []WorkbookOperationable)()
    SetTables(value []WorkbookTableable)()
    SetWorksheets(value []WorkbookWorksheetable)()
}
