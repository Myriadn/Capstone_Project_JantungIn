package dto

type AuthRegisterRequest struct {
	Name        string `json:"name"`
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email"`
	Password    string `json:"password" binding:"required"`
	DateOfBirth string `json:"dateOfBirth"`
}

type AuthLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthLoginEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthVerifyOTPRequest struct {
	UserID string `json:"userId" binding:"required"`
	Code   string `json:"code" binding:"required"`
}

type UpdateProfileRequest struct {
	Name        string `json:"name"`
	DateOfBirth string `json:"dateOfBirth"`
}

type CreateDiagnosisRequest struct {
	PatientID             string  `json:"patientId"` // opsional, hanya untuk admin/dokter
	Age                   int     `json:"age" binding:"required,min=1,max=120"`
	Sex                   string  `json:"sex" binding:"required,oneof=Male Female"`
	ChestPainType         string  `json:"chestPainType" binding:"required"`
	RestingBloodPressure  float64 `json:"restingBloodPressure" binding:"required,gt=0"`
	SerumCholesterol      float64 `json:"serumCholesterol" binding:"required,gte=0"`
	FastingBloodSugar     float64 `json:"fastingBloodSugar" binding:"required,gte=0"`
	RestingEcgResults     string  `json:"restingEcgResults" binding:"required"`
	MaximumHeartRate      int     `json:"maximumHeartRate" binding:"required,gte=0"`
	ExerciseInducedAngina string  `json:"exerciseInducedAngina" binding:"required,oneof=Yes No"`
	StDepression          float64 `json:"stDepression" binding:"gte=0"`
	StSegment             string  `json:"stSegment" binding:"required"`
	MajorVessels          int     `json:"majorVessels" binding:"gte=0,lte=3"`
	Thalassemia           string  `json:"thalassemia" binding:"required"`
}
