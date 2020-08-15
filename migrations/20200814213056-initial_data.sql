-- +migrate Up
INSERT INTO "tag" ("id", "name")
VALUES (1, 'Одежда'),
       (2, 'Деньги'),
       (3, 'Уколы'),
       (4, 'Сходить в магазин'),
       (5, 'Носить сумки'),
       (6, 'Автомощь');

-- +migrate Down
TRUNCATE TABLE "tag" CASCADE;
