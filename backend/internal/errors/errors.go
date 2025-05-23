package errs

import "errors"

var ErrUserAlreadyExists = errors.New("пользователь с такой почтой уже существует")
var ErrEmailValidation = errors.New("email не прошел валидацию")
var ErrWrongPassword = errors.New("неверный пароль")
var ErrPasswordLength = errors.New("длина пароля меньше нужного")

var ErrNotFoundInJWT = errors.New("user_id не найден в токене")
var ErrUnsupportedSignatureMethod = errors.New("неподдерживаемый метод подписи")
var ErrTokenExpired = errors.New("токен истек")
var ErrNoExpClaim = errors.New("не найдено поле exp в токене")
var ErrMoodDescLength = errors.New("длина описания больше допустимого")
var ErrMoodEmotesLength = errors.New("длина эмоций больше допустимого")
var ErrRegistrationDisabled = errors.New("регистрация отключена")
