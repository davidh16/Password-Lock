package service

import "password-lock/models"

func (s Service) FindAllSecurityQuestions() ([]models.PersonalQuestion, error) {
	return s.userRepository.FindAllSecurityQuestions()
}
