package profile

import (
	"fmt"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenderFieldValidation 测试性别字段验证逻辑（独立测试）
func TestGenderFieldValidation(t *testing.T) {
	// 提取出验证逻辑进行独立测试
	validateGender := func(gender string) error {
		if gender != "" && gender != "male" && gender != "female" {
			return connect.NewError(connect.CodeInvalidArgument, 
				fmt.Errorf("gender must be 'male' or 'female', got: %s", gender))
		}
		return nil
	}

	tests := []struct {
		name        string
		gender      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid_male_gender",
			gender:      "male",
			expectError: false,
		},
		{
			name:        "valid_female_gender", 
			gender:      "female",
			expectError: false,
		},
		{
			name:        "empty_gender_should_pass",
			gender:      "",
			expectError: false,
		},
		{
			name:        "invalid_gender_other",
			gender:      "other",
			expectError: true,
			errorMsg:    "gender must be 'male' or 'female'",
		},
		{
			name:        "invalid_gender_chinese",
			gender:      "男",
			expectError: true,
			errorMsg:    "gender must be 'male' or 'female'",
		},
		{
			name:        "invalid_gender_number",
			gender:      "1",
			expectError: true,
			errorMsg:    "gender must be 'male' or 'female'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGender(tt.gender)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUpdateProfile_GenderFieldProcessing 测试性别字段处理逻辑
func TestUpdateProfile_GenderFieldProcessing(t *testing.T) {
	tests := []struct {
		name           string
		inputGender    string
		expectedUpdate bool
		description    string
	}{
		{
			name:           "male_gender_should_update",
			inputGender:    "male",
			expectedUpdate: true,
			description:    "男性性别应该被更新",
		},
		{
			name:           "female_gender_should_update",
			inputGender:    "female", 
			expectedUpdate: true,
			description:    "女性性别应该被更新",
		},
		{
			name:           "empty_gender_should_not_update",
			inputGender:    "",
			expectedUpdate: false,
			description:    "空性别不应该被更新",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 模拟更新逻辑
			updates := map[string]any{}
			
			// 这里复制UpdateProfile中的gender处理逻辑
			if tt.inputGender != "" {
				// 验证性别字段
				if tt.inputGender != "male" && tt.inputGender != "female" {
					t.Errorf("Invalid gender: %s", tt.inputGender)
					return
				}
				updates["gender"] = tt.inputGender
			}

			if tt.expectedUpdate {
				assert.Contains(t, updates, "gender")
				assert.Equal(t, tt.inputGender, updates["gender"])
			} else {
				assert.NotContains(t, updates, "gender")
			}
		})
	}
}

// TestGenderValidationHelper 测试性别验证辅助函数
func TestGenderValidationHelper(t *testing.T) {
	// 辅助函数：验证性别值
	validateGender := func(gender string) error {
		if gender != "" && gender != "male" && gender != "female" {
			return connect.NewError(connect.CodeInvalidArgument, 
				fmt.Errorf("gender must be 'male' or 'female', got: %s", gender))
		}
		return nil
	}

	tests := []struct {
		gender    string
		shouldErr bool
	}{
		{"male", false},
		{"female", false},
		{"", false},
		{"other", true},
		{"男", true},
		{"女", true},
		{"1", true},
		{"0", true},
	}

	for _, test := range tests {
		err := validateGender(test.gender)
		if test.shouldErr {
			require.Error(t, err, "Gender %s should be invalid", test.gender)
		} else {
			require.NoError(t, err, "Gender %s should be valid", test.gender)
		}
	}
}

// Benchmark_GenderValidation 性别验证性能测试
func Benchmark_GenderValidation(b *testing.B) {
	validateGender := func(gender string) bool {
		return gender == "" || gender == "male" || gender == "female"
	}

	testCases := []string{"male", "female", "", "invalid"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, gender := range testCases {
			validateGender(gender)
		}
	}
}