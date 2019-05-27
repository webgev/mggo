package mggo

import (
    "crypto/rand"
    mrand "math/rand"
    "net/smtp"
    "strconv"
    "fmt"
    "time"
)

type validationType string

//  ValidationEmailType email validation
//  ValidationSmsType sms validation
const (
    ValidationEmailType validationType = "email"
    ValidationSmsType   validationType = "sms"
)

// Validation struct
type Validation struct {
    Type  validationType
    Email string
    Phone string
    Token string
    Code  int
}

// Create validation
func (v *Validation) Create() string {
    if v.Token == "" {
        v.GenerateToken()
    }
    v.GenerateCode()
    if v.Type == ValidationEmailType {
        v.sendEmailMessage()
    }
    redicClient.Set(v.Token, string(v.Code), 0)
    redicClient.Expire(v.Token, 10 * time.Minute)
    return v.Token
}

// Verification validation
func (v *Validation) Verification() bool {
    if v.Token == "" || v.Code == 0 {
        return false
    }
    val, err := redicClient.Get(v.Token).Result()
	if err != nil {
		panic(err)
    }
    result := val == string(v.Code) 
    if result {
        redicClient.Del(v.Token)
    }
    return result
}

func (v *Validation) sendEmailMessage() {
    section, err := config.GetSection("smtp")
    if err != nil {
        panic(err)
    }
    email, err := section.GetKey("email")
    if err != nil {
        panic(err)
    }
    pass, err := section.GetKey("password")
    if err != nil {
        panic(err)
    }
    server, err := section.GetKey("server")
    if err != nil {
        panic(err)
    }
    port, err := section.GetKey("port")
    if err != nil {
        panic(err)
    }
    from := email.String()
    body := v.Code
    to := v.Email
    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: Validation\n\n" +
        strconv.Itoa(body)

    err = smtp.SendMail(server.String() + ":" + port.String(),
        smtp.PlainAuth("", from, pass.String(), server.String()),
        from, []string{to}, []byte(msg))

    if err != nil {
        panic(ErrorInternalServer{Message: err.Error()})
    }
    // TODO: Записать в базу или редис
}

// GenerateToken is generate token for validation
func (v *Validation) GenerateToken() {
    b := make([]byte, 30)
	rand.Read(b)
    v.Token = fmt.Sprintf("valid-%x", b)
    
}

// GenerateCode is generate gode for validation
func (v *Validation) GenerateCode() {
    mrand.Seed(time.Now().UnixNano())
    min := 100000
    max := 999999
    v.Code = mrand.Intn(max - min) + min
}
