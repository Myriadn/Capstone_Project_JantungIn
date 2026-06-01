package dto

// PatientResponse digunakan untuk endpoint GET /api/v1/admin/patients
type PatientResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       *string `json:"email,omitempty"`
	DateOfBirth *string `json:"dateOfBirth,omitempty"`
	Role        string  `json:"role"`
	CreatedAt   string  `json:"createdAt"`
}

type AuthUserResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email"`
	Role     string  `json:"role"`
}

type AuthRegisterResponse struct {
	StatusCode int              `json:"statusCode"`
	Message    string           `json:"message"`
	Data       AuthRegisterData `json:"data"`
}

type AuthRegisterData struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email"`
	Role     string  `json:"role"`
}

type AuthLoginResponse struct {
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Data       AuthLoginData `json:"data"`
}

type AuthLoginData struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email"`
	Role     string  `json:"role"`
	Token    string  `json:"token"`
}

type AuthLoginOTPData struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Username     *string `json:"username,omitempty"`
	Email        *string `json:"email"`
	Role         string  `json:"role"`
	OTPRequired  bool    `json:"otpRequired"`
	OTPExpiresIn int     `json:"otpExpiresIn"`
}

type ProfileResponse struct {
	StatusCode int              `json:"statusCode"`
	Message    string           `json:"message"`
	Data       AuthUserResponse `json:"data"`
}

type UpdateProfileResponse struct {
	StatusCode int               `json:"statusCode"`
	Message    string            `json:"message"`
	Data       UpdateProfileData `json:"data"`
}

type UpdateProfileData struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Username    *string `json:"username,omitempty"`
	Email       *string `json:"email"`
	DateOfBirth *string `json:"dateOfBirth,omitempty"`
}

type DiagnosisResultData struct {
	ID                 string  `json:"id"`
	UserID             string  `json:"userId"`
	ResultPercentage   float64 `json:"resultPercentage"`
	CardiovascularRisk string  `json:"cardiovascularRisk"`
	Prediction         string  `json:"prediction"`
	CreatedAt          string  `json:"createdAt"`
}

type DiagnosisUserInfo struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type DiagnosisResponse struct {
	ID                    string             `json:"id"`
	UserID                string             `json:"userId"`
	CreatedBy             *string            `json:"createdBy"`
	Age                   int                `json:"age"`
	Sex                   string             `json:"sex"`
	ChestPainType         string             `json:"chestPainType"`
	RestingEcgResults     string             `json:"restingEcgResults"`
	FastingBloodSugar     float64            `json:"fastingBloodSugar"`
	RestingBloodPressure  float64            `json:"restingBloodPressure"`
	MaximumHeartRate      int                `json:"maximumHeartRate"`
	ExerciseInducedAngina string             `json:"exerciseInducedAngina"`
	StSegment             string             `json:"stSegment"`
	MajorVessels          int                `json:"majorVessels"`
	Thalassemia           string             `json:"thalassemia"`
	SerumCholesterol      float64            `json:"serumCholesterol"`
	StDepression          float64            `json:"stDepression"`
	ResultPercentage      float64            `json:"resultPercentage"`
	CardiovascularRisk    string             `json:"cardiovascularRisk"`
	Prediction            string             `json:"prediction"`
	Patient               DiagnosisUserInfo  `json:"patient"`
	Creator               *DiagnosisUserInfo `json:"creator,omitempty"`
	CreatedAt             string             `json:"createdAt"`
	UpdatedAt             string             `json:"updatedAt"`
}
