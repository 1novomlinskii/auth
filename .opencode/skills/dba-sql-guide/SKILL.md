---
name: dba-sql-guide
description: Форматирование SQL по руководству стиля sqlstyle.guide — ключевые слова UPPERCASE, snake_case, коридорное выравнивание, правильные отступы
license: CC-BY-SA-4.0
compatibility: opencode
---

# SQL Style Guide — инструкция форматирования

## Регистр

- Все зарезервированные ключевые слова (`SELECT`, `FROM`, `WHERE`, `JOIN`, `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `ALTER`, `DROP`, `GROUP BY`, `ORDER BY`, `HAVING`, `SET`, `VALUES`, `INNER`, `LEFT`, `RIGHT`, `FULL`, `ON`, `AND`, `OR`, `IN`, `BETWEEN`, `LIKE`, `CASE`, `WHEN`, `THEN`, `ELSE`, `END`, `AS`, `NOT`, `NULL`, `IS`, `DISTINCT`, `UNION`, `ALL`, `EXISTS`, `LIMIT`, `OFFSET`, `RETURNING` и т.д.) — **только UPPERCASE**.
- Идентификаторы (таблицы, столбцы, алиасы) — **только lowercase**, `snake_case`.
- camelCase запрещён.

## Выравнивание (коридор)

Главные ключевые слова выравниваются по правому краю, всё остальное — по левому:
```
SELECT column_a,
       column_b
  FROM table_a AS a
 WHERE a.column_a = 1
   AND a.column_b = 2;
```

## Пробелы и пустые строки

- Пробелы до и после `=`
- Пробел после запятой `,`
- `AND` / `OR` — всегда с новой строки, с отступом
- После каждой основной клаузы — перенос строки
- Пустая строка между логическими блоками (например, между `JOIN`)

## Отступы

- В `CREATE TABLE` — 4 пробела для столбцов
- Объявления столбцов и ограничений выравниваются колонками:
  ```
  CREATE TABLE staff (
      PRIMARY KEY (staff_num),
      staff_num      INT(5)       NOT NULL,
      first_name     VARCHAR(100) NOT NULL
  );
  ```

## Псевдонимы

- Всегда с `AS`
- Короткие: из первых букв слов имени объекта, с цифрой при дублировании
  ```
  SELECT first_name AS fn
    FROM staff AS s1;
  ```

## JOIN

- `JOIN` — справа от коридора, каждый с новой строки
- Дополнительные условия `ON` — с новой строки с отступом
  ```
  SELECT r.last_name
    FROM riders AS r
         INNER JOIN bikes AS b
         ON r.bike_vin_num = b.vin_num
            AND b.engine_tally > 2;
  ```

## Подзапросы

- Выравниваются по правому краю коридора
- Внутри — те же правила
- Закрывающая скобка — на уровне открывающей:
  ```
  WHERE r.last_name IN
        (SELECT c.last_name
           FROM champions AS c
          WHERE c.confirmed = 'Y');
  ```

## Конструкции

- `BETWEEN` вместо цепочек `AND`
- `IN()` вместо цепочек `OR`
- `CASE` для интерпретации значений
- `UNION ALL` (предпочитать обычному `UNION`, если дубликаты не мешают)

## Имена таблиц

- Собирательные существительные (`staff`, `crew`), без префиксов (`tbl_`)
- Без совпадения с названиями столбцов

## Имена столбцов

- Единственное число
- Суффиксы: `_id`, `_status`, `_total`, `_num`, `_name`, `_seq`, `_date`, `_tally`, `_size`, `_addr`
- Не использовать `id` как имя первичного ключа

## Хранимые процедуры

- Имя содержит глагол
- Без префикса `sp_`

## CREATE TABLE — порядок

1. `PRIMARY KEY` — сразу после `CREATE TABLE`
2. Столбцы с типами, `DEFAULT`, `NOT NULL`
3. Ограничения — под столбцом, к которому относятся
4. Табличные ограничения — в конце

## Типы данных

- Стандартные ANSI SQL (не vendor-specific)
- `REAL` / `FLOAT` только для плавающей точки, иначе `NUMERIC` / `DECIMAL`

## Даты

- ISO 8601: `YYYY-MM-DDTHH:MM:SS.SSSSS`

## Комментарии

- `/* */` (C-style) или `--` (до конца строки)

## Чего избегать

- `camelCase` — нечитаем
- Префиксы `sp_`, `tbl_` — избыточны
- EAV (Entity-Attribute-Value)
- Принципы ООП в схемах БД
- Разделение логически одной таблицы по времени или географии
- Vendor-specific ключевые слова, если есть ANSI-аналог
