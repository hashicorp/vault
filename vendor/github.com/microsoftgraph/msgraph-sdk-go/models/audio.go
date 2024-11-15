package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e "github.com/microsoft/kiota-abstractions-go/store"
)

type Audio struct {
    // Stores model information.
    backingStore ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore
}
// NewAudio instantiates a new Audio and sets the default values.
func NewAudio()(*Audio) {
    m := &Audio{
    }
    m.backingStore = ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStoreFactoryInstance();
    m.SetAdditionalData(make(map[string]any))
    return m
}
// CreateAudioFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
// returns a Parsable when successful
func CreateAudioFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewAudio(), nil
}
// GetAdditionalData gets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
// returns a map[string]any when successful
func (m *Audio) GetAdditionalData()(map[string]any) {
    val , err :=  m.backingStore.Get("additionalData")
    if err != nil {
        panic(err)
    }
    if val == nil {
        var value = make(map[string]any);
        m.SetAdditionalData(value);
    }
    return val.(map[string]any)
}
// GetAlbum gets the album property value. The title of the album for this audio file.
// returns a *string when successful
func (m *Audio) GetAlbum()(*string) {
    val, err := m.GetBackingStore().Get("album")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetAlbumArtist gets the albumArtist property value. The artist named on the album for the audio file.
// returns a *string when successful
func (m *Audio) GetAlbumArtist()(*string) {
    val, err := m.GetBackingStore().Get("albumArtist")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetArtist gets the artist property value. The performing artist for the audio file.
// returns a *string when successful
func (m *Audio) GetArtist()(*string) {
    val, err := m.GetBackingStore().Get("artist")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetBackingStore gets the BackingStore property value. Stores model information.
// returns a BackingStore when successful
func (m *Audio) GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore) {
    return m.backingStore
}
// GetBitrate gets the bitrate property value. Bitrate expressed in kbps.
// returns a *int64 when successful
func (m *Audio) GetBitrate()(*int64) {
    val, err := m.GetBackingStore().Get("bitrate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetComposers gets the composers property value. The name of the composer of the audio file.
// returns a *string when successful
func (m *Audio) GetComposers()(*string) {
    val, err := m.GetBackingStore().Get("composers")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetCopyright gets the copyright property value. Copyright information for the audio file.
// returns a *string when successful
func (m *Audio) GetCopyright()(*string) {
    val, err := m.GetBackingStore().Get("copyright")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetDisc gets the disc property value. The number of the disc this audio file came from.
// returns a *int32 when successful
func (m *Audio) GetDisc()(*int32) {
    val, err := m.GetBackingStore().Get("disc")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDiscCount gets the discCount property value. The total number of discs in this album.
// returns a *int32 when successful
func (m *Audio) GetDiscCount()(*int32) {
    val, err := m.GetBackingStore().Get("discCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetDuration gets the duration property value. Duration of the audio file, expressed in milliseconds
// returns a *int64 when successful
func (m *Audio) GetDuration()(*int64) {
    val, err := m.GetBackingStore().Get("duration")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int64)
    }
    return nil
}
// GetFieldDeserializers the deserialization information for the current model
// returns a map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error) when successful
func (m *Audio) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["album"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAlbum(val)
        }
        return nil
    }
    res["albumArtist"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetAlbumArtist(val)
        }
        return nil
    }
    res["artist"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetArtist(val)
        }
        return nil
    }
    res["bitrate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetBitrate(val)
        }
        return nil
    }
    res["composers"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetComposers(val)
        }
        return nil
    }
    res["copyright"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetCopyright(val)
        }
        return nil
    }
    res["disc"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDisc(val)
        }
        return nil
    }
    res["discCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDiscCount(val)
        }
        return nil
    }
    res["duration"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt64Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetDuration(val)
        }
        return nil
    }
    res["genre"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetGenre(val)
        }
        return nil
    }
    res["hasDrm"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetHasDrm(val)
        }
        return nil
    }
    res["isVariableBitrate"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetIsVariableBitrate(val)
        }
        return nil
    }
    res["@odata.type"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetOdataType(val)
        }
        return nil
    }
    res["title"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetStringValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTitle(val)
        }
        return nil
    }
    res["track"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrack(val)
        }
        return nil
    }
    res["trackCount"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetTrackCount(val)
        }
        return nil
    }
    res["year"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetInt32Value()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetYear(val)
        }
        return nil
    }
    return res
}
// GetGenre gets the genre property value. The genre of this audio file.
// returns a *string when successful
func (m *Audio) GetGenre()(*string) {
    val, err := m.GetBackingStore().Get("genre")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetHasDrm gets the hasDrm property value. Indicates if the file is protected with digital rights management.
// returns a *bool when successful
func (m *Audio) GetHasDrm()(*bool) {
    val, err := m.GetBackingStore().Get("hasDrm")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetIsVariableBitrate gets the isVariableBitrate property value. Indicates if the file is encoded with a variable bitrate.
// returns a *bool when successful
func (m *Audio) GetIsVariableBitrate()(*bool) {
    val, err := m.GetBackingStore().Get("isVariableBitrate")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*bool)
    }
    return nil
}
// GetOdataType gets the @odata.type property value. The OdataType property
// returns a *string when successful
func (m *Audio) GetOdataType()(*string) {
    val, err := m.GetBackingStore().Get("odataType")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTitle gets the title property value. The title of the audio file.
// returns a *string when successful
func (m *Audio) GetTitle()(*string) {
    val, err := m.GetBackingStore().Get("title")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*string)
    }
    return nil
}
// GetTrack gets the track property value. The number of the track on the original disc for this audio file.
// returns a *int32 when successful
func (m *Audio) GetTrack()(*int32) {
    val, err := m.GetBackingStore().Get("track")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetTrackCount gets the trackCount property value. The total number of tracks on the original disc for this audio file.
// returns a *int32 when successful
func (m *Audio) GetTrackCount()(*int32) {
    val, err := m.GetBackingStore().Get("trackCount")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// GetYear gets the year property value. The year the audio file was recorded.
// returns a *int32 when successful
func (m *Audio) GetYear()(*int32) {
    val, err := m.GetBackingStore().Get("year")
    if err != nil {
        panic(err)
    }
    if val != nil {
        return val.(*int32)
    }
    return nil
}
// Serialize serializes information the current object
func (m *Audio) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteStringValue("album", m.GetAlbum())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("albumArtist", m.GetAlbumArtist())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("artist", m.GetArtist())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("bitrate", m.GetBitrate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("composers", m.GetComposers())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("copyright", m.GetCopyright())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("disc", m.GetDisc())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("discCount", m.GetDiscCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt64Value("duration", m.GetDuration())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("genre", m.GetGenre())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("hasDrm", m.GetHasDrm())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteBoolValue("isVariableBitrate", m.GetIsVariableBitrate())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("@odata.type", m.GetOdataType())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteStringValue("title", m.GetTitle())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("track", m.GetTrack())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("trackCount", m.GetTrackCount())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteInt32Value("year", m.GetYear())
        if err != nil {
            return err
        }
    }
    {
        err := writer.WriteAdditionalData(m.GetAdditionalData())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the AdditionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *Audio) SetAdditionalData(value map[string]any)() {
    err := m.GetBackingStore().Set("additionalData", value)
    if err != nil {
        panic(err)
    }
}
// SetAlbum sets the album property value. The title of the album for this audio file.
func (m *Audio) SetAlbum(value *string)() {
    err := m.GetBackingStore().Set("album", value)
    if err != nil {
        panic(err)
    }
}
// SetAlbumArtist sets the albumArtist property value. The artist named on the album for the audio file.
func (m *Audio) SetAlbumArtist(value *string)() {
    err := m.GetBackingStore().Set("albumArtist", value)
    if err != nil {
        panic(err)
    }
}
// SetArtist sets the artist property value. The performing artist for the audio file.
func (m *Audio) SetArtist(value *string)() {
    err := m.GetBackingStore().Set("artist", value)
    if err != nil {
        panic(err)
    }
}
// SetBackingStore sets the BackingStore property value. Stores model information.
func (m *Audio) SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)() {
    m.backingStore = value
}
// SetBitrate sets the bitrate property value. Bitrate expressed in kbps.
func (m *Audio) SetBitrate(value *int64)() {
    err := m.GetBackingStore().Set("bitrate", value)
    if err != nil {
        panic(err)
    }
}
// SetComposers sets the composers property value. The name of the composer of the audio file.
func (m *Audio) SetComposers(value *string)() {
    err := m.GetBackingStore().Set("composers", value)
    if err != nil {
        panic(err)
    }
}
// SetCopyright sets the copyright property value. Copyright information for the audio file.
func (m *Audio) SetCopyright(value *string)() {
    err := m.GetBackingStore().Set("copyright", value)
    if err != nil {
        panic(err)
    }
}
// SetDisc sets the disc property value. The number of the disc this audio file came from.
func (m *Audio) SetDisc(value *int32)() {
    err := m.GetBackingStore().Set("disc", value)
    if err != nil {
        panic(err)
    }
}
// SetDiscCount sets the discCount property value. The total number of discs in this album.
func (m *Audio) SetDiscCount(value *int32)() {
    err := m.GetBackingStore().Set("discCount", value)
    if err != nil {
        panic(err)
    }
}
// SetDuration sets the duration property value. Duration of the audio file, expressed in milliseconds
func (m *Audio) SetDuration(value *int64)() {
    err := m.GetBackingStore().Set("duration", value)
    if err != nil {
        panic(err)
    }
}
// SetGenre sets the genre property value. The genre of this audio file.
func (m *Audio) SetGenre(value *string)() {
    err := m.GetBackingStore().Set("genre", value)
    if err != nil {
        panic(err)
    }
}
// SetHasDrm sets the hasDrm property value. Indicates if the file is protected with digital rights management.
func (m *Audio) SetHasDrm(value *bool)() {
    err := m.GetBackingStore().Set("hasDrm", value)
    if err != nil {
        panic(err)
    }
}
// SetIsVariableBitrate sets the isVariableBitrate property value. Indicates if the file is encoded with a variable bitrate.
func (m *Audio) SetIsVariableBitrate(value *bool)() {
    err := m.GetBackingStore().Set("isVariableBitrate", value)
    if err != nil {
        panic(err)
    }
}
// SetOdataType sets the @odata.type property value. The OdataType property
func (m *Audio) SetOdataType(value *string)() {
    err := m.GetBackingStore().Set("odataType", value)
    if err != nil {
        panic(err)
    }
}
// SetTitle sets the title property value. The title of the audio file.
func (m *Audio) SetTitle(value *string)() {
    err := m.GetBackingStore().Set("title", value)
    if err != nil {
        panic(err)
    }
}
// SetTrack sets the track property value. The number of the track on the original disc for this audio file.
func (m *Audio) SetTrack(value *int32)() {
    err := m.GetBackingStore().Set("track", value)
    if err != nil {
        panic(err)
    }
}
// SetTrackCount sets the trackCount property value. The total number of tracks on the original disc for this audio file.
func (m *Audio) SetTrackCount(value *int32)() {
    err := m.GetBackingStore().Set("trackCount", value)
    if err != nil {
        panic(err)
    }
}
// SetYear sets the year property value. The year the audio file was recorded.
func (m *Audio) SetYear(value *int32)() {
    err := m.GetBackingStore().Set("year", value)
    if err != nil {
        panic(err)
    }
}
type Audioable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackedModel
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAlbum()(*string)
    GetAlbumArtist()(*string)
    GetArtist()(*string)
    GetBackingStore()(ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)
    GetBitrate()(*int64)
    GetComposers()(*string)
    GetCopyright()(*string)
    GetDisc()(*int32)
    GetDiscCount()(*int32)
    GetDuration()(*int64)
    GetGenre()(*string)
    GetHasDrm()(*bool)
    GetIsVariableBitrate()(*bool)
    GetOdataType()(*string)
    GetTitle()(*string)
    GetTrack()(*int32)
    GetTrackCount()(*int32)
    GetYear()(*int32)
    SetAlbum(value *string)()
    SetAlbumArtist(value *string)()
    SetArtist(value *string)()
    SetBackingStore(value ie8677ce2c7e1b4c22e9c3827ecd078d41185424dd9eeb92b7d971ed2d49a392e.BackingStore)()
    SetBitrate(value *int64)()
    SetComposers(value *string)()
    SetCopyright(value *string)()
    SetDisc(value *int32)()
    SetDiscCount(value *int32)()
    SetDuration(value *int64)()
    SetGenre(value *string)()
    SetHasDrm(value *bool)()
    SetIsVariableBitrate(value *bool)()
    SetOdataType(value *string)()
    SetTitle(value *string)()
    SetTrack(value *int32)()
    SetTrackCount(value *int32)()
    SetYear(value *int32)()
}
