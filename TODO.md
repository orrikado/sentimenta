# backend

- [ ] api смены пароля который берёт(отдельный от /api/user/update)

```json
{
  currentPassword = "12345678",
  newPassword = "87654321"
}
```

- [ ] /api/user/update должен требовать пароль ЕСЛИ юзер пытается изменить почту
- [ ] замержить все /api/user/{get, update} в один и просто юзать разные http запросы по одному пути /api/user
  - [ ] /api/user/update должен использовать patch а не put, пут для новых данных а не изменения
- [ ] добавить констрейнт для юзера что нельзя несколько людей с одним ником
- [ ] добавить полноценный реджекс для проверки имейла на бекэнд [regex](https://regex101.com/r/6EL6YF/1)

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
