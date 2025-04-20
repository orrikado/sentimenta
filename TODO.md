# backend

- [x] api смены пароля который берёт(отдельный от /api/user/update)

```json
{
  password = "12345678",
  new_password = "87654321"
}
```

- [x] /api/user/update должен требовать пароль ЕСЛИ юзер пытается изменить почту
- [x] /api/user/update должен использовать patch а не put, пут для новых данных а не изменения
- [x] добавить констрейнт для юзера что нельзя несколько людей с одним ником
- [x] добавить полноценный реджекс для проверки имейла на бекэнд [regex](https://regex101.com/r/6EL6YF/1)

# frontend

- [ ] gravitar
- [ ] graph type shi

# ci/cd

- [ ] написать ci/cd pipeline
  - [ ] dependabot
  - [ ] build docker image
  - [ ] lint
    - [ ] front
    - [ ] back
