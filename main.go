package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
)

func main() {
	port := os.Getenv("PORT")
	router := mux.NewRouter()
	secure := router.PathPrefix("/secure").Subrouter()
	secure.Use(JwtVerify)
	secure.HandleFunc("/tokenHook", processTokenInlineHook).Methods("POST")
	http.ListenAndServe(":"+port, router)
}

func processTokenInlineHook(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	fmt.Println("processTokenInlineHook : Token inline Hook invoked..")
	w.Header().Set("Content-Type", "application/json")

	var tp TokenPayload

	err := json.NewDecoder(r.Body).Decode(&tp)
	if err != nil {
		fmt.Println("Error parsing httpBody", err)
	}

	fmt.Printf("TP value is %+v \n", tp)

	json.Marshal(tp)

	audClaim := strings.Split(tp.Data.Access.Claims.Aud, ",")

	commandValue_1 := CommandValue{
		Op:    "add",
		Path:  "/claims/aud1",
		Value: fmt.Sprintf("%s1", audClaim[0]),
	}

	u, err := url.Parse(tp.Data.Context.Request.URL.Value)
	if err != nil {
		log.Fatal(err)
	}

	queryParams := u.Query()

	fmt.Println("Resource URL Query Parameters : ", queryParams)

	command_1 := Command{
		CommandType: "com.okta.access.patch",
		Vaue:        []CommandValue{commandValue_1},
	}

	isError := false
	errorSummary := ""

	var commands []Command

	appendCommands(&commands, &command_1)

	if isError {
		commandValue_err := CommandValue{
			Op:    "add",
			Path:  "/claims/isError",
			Value: "true",
		}
		command_err := Command{
			CommandType: "com.okta.access.patch",
			Vaue:        []CommandValue{commandValue_err},
		}
		appendCommands(&commands, &command_err)
	}

	response := Response{
		Commands: commands,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	fmt.Println("Response to Okta = ", string(jsonResponse))

	if isError {
		fmt.Println("**** Error *** : " + errorSummary)
		w.Write([]byte(fmt.Sprintf("{\"error\":{\"errorSummary\":\"%s\"}}", errorSummary)))
	}

	if !isError {
		fmt.Println("**** Responding ***")
		w.Write(jsonResponse)
		fmt.Println("**** Response Successful ***")
	}

	end := time.Now()

	nanoTimeDelta := end.UnixNano() - start.UnixNano()
	millisDelta := nanoTimeDelta / 1000000
	fmt.Printf("Total [%d] nano(s), [%d] milli(s) taken to complete token enrichment\n", nanoTimeDelta, millisDelta)
}

func appendCommands(commands *[]Command, command *Command) {
	*commands = append(*commands, *command)
}

type Response struct {
	Commands []Command `json:"commands"`
}

type Command struct {
	CommandType string         `json:"type"`
	Vaue        []CommandValue `json:"value"`
}

type CommandValue struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ParseHttpHeader(r)

		access_token := r.Header.Get("Authorization")

		if len(access_token) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("AuthZ Error : Need an access token to use this service")
			return
		}

		if strings.HasPrefix(access_token, "Bearer ") {
			access_token = strings.TrimPrefix(access_token, "Bearer ")
			fmt.Println("Access Token is : ", access_token)
		}

		toValidate := map[string]string{}
		toValidate["aud"] = os.Getenv("JWT_AT_AUD")
		toValidate["cid"] = os.Getenv("JWT_AT_CLIENT_ID")
		toValidate["sub"] = os.Getenv("JWT_AT_CLIENT_ID")

		jwtVerifierSetup := jwtverifier.JwtVerifier{
			Issuer:           os.Getenv("JWT_AT_ISS"),
			ClaimsToValidate: toValidate,
		}

		verifier := jwtVerifierSetup.New()
		verifier.SetLeeway("2m")

		token, err := verifier.VerifyAccessToken(access_token)

		if err != nil {
			fmt.Println("Error validating access token : ", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(fmt.Sprintf("AuthZ Error : Invalid access token : %v", err))
			return
		}
		scopes := token.Claims["scp"].([]interface{})

		requiredScope := os.Getenv("JWT_AT_REQ_SCOPE")
		var hasScope bool
		for _, scope := range scopes {
			if scope == requiredScope {
				hasScope = true
				fmt.Println("Has required scope", scope)
				break
			}
		}

		if hasScope {
			fmt.Println("Invoke token entichment...")
			next.ServeHTTP(w, r)
			fmt.Println("Done Token enrichment.")
		} else {
			fmt.Println("Invalid access token...")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Required scopes not found")
			fmt.Println("Aboth token enrichment...")
		}

	})
}

func ParseHttpBody(r *http.Request) {
	jsonMap := make(map[string]interface{})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&jsonMap)
	if err != nil {
		panic(err)
	}
	fmt.Println("Payload from Okta : ", jsonMap)
}

func ParseHttpHeader(r *http.Request) {

	fmt.Println("***Headers [Start]***")
	headers := r.Header

	for k, v := range headers {
		fmt.Println(k, " : ", v)
	}
	fmt.Println("***Headers [End]***")

	fmt.Println("***Cookie [Start]***")
	cookies := r.Cookies()

	for cookie := range cookies {
		fmt.Println("Cookie : ", cookie)
	}
	fmt.Println("***Cookie [End]***")

	fmt.Println("***Content Length : ", r.ContentLength)

}

type TokenPayload struct {
	Source            string      `json:"source,omitempty"`
	EventID           string      `json:"eventId,omitempty"`
	EventTime         time.Time   `json:"eventTime,omitempty"`
	EventTypeVersion  string      `json:"eventTypeVersion,omitempty"`
	CloudEventVersion string      `json:"cloudEventVersion,omitempty"`
	ContentType       string      `json:"contentType,omitempty"`
	EventType         string      `json:"eventType,omitempty"`
	Data              DataElement `json:"data,omitempty"`
}

type DataElement struct {
	Identity IdentityElement `json:"identity,omitempty"`
	Access   AccessElement   `json:"access,omitempty"`
	Context  ContextElement  `json:"context,omitempty"`
}

type IdentityElement struct {
	Claims struct {
		Sub               string   `json:"sub,omitempty"`
		Name              string   `json:"name,omitempty"`
		Email             string   `json:"email,omitempty"`
		Ver               int      `json:"ver,omitempty"`
		Iss               string   `json:"iss,omitempty"`
		Aud               string   `json:"aud,omitempty"`
		Jti               string   `json:"jti,omitempty"`
		Amr               []string `json:"amr,omitempty"`
		Idp               string   `json:"idp,omitempty"`
		Nonce             string   `json:"nonce,omitempty"`
		PreferredUsername string   `json:"preferred_username,omitempty"`
		AuthTime          int      `json:"auth_time,omitempty"`
	} `json:"claims,omitempty"`
	Token struct {
		Lifetime struct {
			Expiration int `json:"expiration,omitempty"`
		} `json:"lifetime,omitempty"`
	} `json:"token,omitempty"`
}

type AccessElement struct {
	Claims struct {
		Ver               int    `json:"ver,omitempty"`
		Jti               string `json:"jti,omitempty"`
		Iss               string `json:"iss,omitempty"`
		Aud               string `json:"aud,omitempty"`
		Cid               string `json:"cid,omitempty"`
		UID               string `json:"uid,omitempty"`
		Sub               string `json:"sub,omitempty"`
		FirstName         string `json:"firstName,omitempty"`
		PreferredUsername string `json:"preferred_username,omitempty"`
		Scope             string `json:"scope,omitempty"`
	} `json:"claims,omitempty"`
	Token struct {
		Lifetime struct {
			Expiration int `json:"expiration,omitempty"`
		} `json:"lifetime,omitempty"`
	}
}

type ContextElement struct {
	Request  RequestElement  `json:"request,omitempty"`
	User     UserElement     `json:"user,omitempty"`
	Protocol ProtocolElement `json:"protocol,omitempty"`
	Policy   PolicyElement   `json:"policy,omitempty"`
	Session  SessionElement  `json:"session,omitempty"`
}

type RequestElement struct {
	ID     string `json:"id"`
	Method string `json:"method,omitempty"`
	URL    struct {
		Value string `json:"value,omitempty"`
	} `json:"url,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
}

type UserElement struct {
	ID              string         `json:"id"`
	PasswordChanged time.Time      `json:"passwordChanged,omitempty"`
	Profile         ProfileElement `json:"profile,omitempty"`
}

type ProtocolElement struct {
	Type    string        `json:"type,omitempty"`
	Client  ClientElement `json:"client,omitempty"`
	Issuer  IssuerElement `json:"issuer,omitempty"`
	Request struct {
		Scope        string `json:"scope,omitempty"`
		State        string `json:"state,omitempty"`
		RedirectURI  string `json:"redirect_uri,omitempty"`
		ResponseMode string `json:"response_mode,omitempty"`
		ResponseType string `json:"response_type,omitempty"`
		ClientID     string `json:"client_id,omitempty"`
	} `json:"request,omitempty"`
}

type IssuerElement struct {
	URI string `json:"uri,omitempty"`
}

type ClientElement struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type ProfileElement struct {
	Login     string `json:"login,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Locale    string `json:"locale,omitempty"`
	TimeZone  string `json:"timeZone,omitempty"`
}

type PolicyElement struct {
	ID   string      `json:"id,omitempty"`
	Rule RuleElement `json:"rule,omitempty"`
}

type RuleElement struct {
	ID string `json:"id,omitempty"`
}

type SessionElement struct {
	ID                       string     `json:"id,omitempty"`
	UserID                   string     `json:"userId,omitempty"`
	Login                    string     `json:"login,omitempty"`
	CreatedAt                time.Time  `json:"createdAt,omitempty"`
	ExpiresAt                time.Time  `json:"expiresAt,omitempty"`
	Status                   string     `json:"status,omitempty"`
	LastPasswordVerification time.Time  `json:"lastPasswordVerification,omitempty"`
	Amr                      []string   `json:"amr,omitempty"`
	Idp                      IdpElement `json:"idp,omitempty"`
	MfaActive                bool       `json:"mfaActive,omitempty"`
}

type IdpElement struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type Error struct {
	ErrorSummary string       `json:"errorSummary"`
	ErrorCauses  []ErrorCause `json:"errorCauses"`
}

type ErrorCause struct {
	ErrorSummary string `json:"errorSummary"`
	Reason       string `json:"reason"`
	LocationType string `json:"locationType"`
	Location     string `json:"location"`
	Domain       string `json:"domain"`
}
