package errs

import "errors"

var ErrUserAlreadyExists = errors.New("пользователь с такой почтой уже существует")
var ErrEmailValidation = errors.New("email не прошел валидацию")
var ErrWrongPassword = errors.New("неверный пароль")

var ErrNotFoundInJWT = errors.New("user_id не найден в токене")
var ErrUnsupportedSignatureMethod = errors.New("неподдерживаемый метод подписи")
