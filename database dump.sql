DROP TABLE IF EXISTS "competence";

DROP SEQUENCE IF EXISTS competence_id_seq;

CREATE SEQUENCE competence_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."competence" (
    "id" integer DEFAULT nextval('competence_id_seq') NOT NULL,
    "name" character(250) NOT NULL,
    "chort_name" character(100) NOT NULL,
    CONSTRAINT "competence_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO
    "competence" ("id", "name", "chort_name")
VALUES
    (1, 'управление проектами', 'уп'),
    (2, 'работа в команде', 'рвк'),
    (3, 'системное и критическое мышление', 'сикм'),
    (4, 'самоорганизация и саморазвитие', 'сис');

DROP TABLE IF EXISTS "message";

DROP SEQUENCE IF EXISTS "Message_id_seq";

CREATE SEQUENCE "Message_id_seq" INCREMENT MINVALUE MAXVALUE CACHE;

CREATE TABLE "public"."message" (
    "id" integer DEFAULT nextval('"Message_id_seq"') NOT NULL,
    "sender" integer NOT NULL,
    "recipient" integer NOT NULL,
    "date" date NOT NULL,
    "time" time without time zone NOT NULL,
    "is_read" boolean DEFAULT false NOT NULL,
    "text_id" integer,
    CONSTRAINT "message_id" PRIMARY KEY ("id")
) WITH (oids = false);

DROP TABLE IF EXISTS "message_text";

DROP SEQUENCE IF EXISTS "Message_text_id_seq";

CREATE SEQUENCE "Message_text_id_seq" INCREMENT MINVALUE MAXVALUE CACHE;

CREATE TABLE "public"."message_text" (
    "id" integer DEFAULT nextval('"Message_text_id_seq"') NOT NULL,
    "text" character(255) NOT NULL,
    "next" integer,
    CONSTRAINT "message_text_id" PRIMARY KEY ("id")
) WITH (oids = false);

DROP TABLE IF EXISTS "sessions";

CREATE TABLE "public"."sessions" (
    "user_id" integer NOT NULL,
    "name" character(250) NOT NULL,
    "secret_key" character(32) NOT NULL,
    "creation_date" date NOT NULL,
    "last_date" date NOT NULL,
    "last_time" time without time zone NOT NULL,
    CONSTRAINT "sessions_secret_key" PRIMARY KEY ("secret_key")
) WITH (oids = false);

DROP TABLE IF EXISTS "tasks";

DROP SEQUENCE IF EXISTS tasks_id_seq;

CREATE SEQUENCE tasks_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."tasks" (
    "id" integer DEFAULT nextval('tasks_id_seq') NOT NULL,
    "name" character(200) NOT NULL,
    "description" character(255) NOT NULL,
    "competence_id" integer NOT NULL,
    "html_num" integer NOT NULL,
    "mark" integer NOT NULL,
    CONSTRAINT "tasks_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

INSERT INTO
    "tasks" (
        "id",
        "name",
        "description",
        "competence_id",
        "html_num",
        "mark"
    )
VALUES
    (
        1,
        'Соотнести сотрудников с soft-навыками',
        'Необходимо соотнести, насколько каждый из сотрудников подойдёт под ту или иную роль.',
        1,
        1,
        2
    ),
    (
        2,
        'Упорядочивание этапов проекта',
        'Вам даны этапы проекта. Необходимо поставить их в верном порядке.',
        1,
        2,
        2
    );

DROP TABLE IF EXISTS "users";

DROP SEQUENCE IF EXISTS "Users_id_seq";

CREATE SEQUENCE "Users_id_seq" INCREMENT MINVALUE MAXVALUE CACHE;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('"Users_id_seq"') NOT NULL,
    "login" character(50) NOT NULL,
    "password" character(250) NOT NULL,
    "full_name" character(200) NOT NULL,
    "email" character(200) NOT NULL,
    "faculty" character(10) DEFAULT '' NOT NULL,
    "study_group" character(10) DEFAULT '' NOT NULL,
    "avatar_id" integer DEFAULT '0' NOT NULL,
    "registration_date" date NOT NULL,
    CONSTRAINT "users_id" PRIMARY KEY ("id"),
    CONSTRAINT "users_login" UNIQUE ("login")
) WITH (oids = false);

INSERT INTO
    "users" (
        "id",
        "login",
        "password",
        "full_name",
        "email",
        "faculty",
        "study_group",
        "avatar_id",
        "registration_date"
    )
VALUES
    (0, '', '', '', '', '', '', 0, '1970-01-01');

ALTER TABLE
    ONLY "public"."message"
ADD
    CONSTRAINT "message_text_id_fkey" FOREIGN KEY (text_id) REFERENCES message_text(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;

ALTER TABLE
    ONLY "public"."message_text"
ADD
    CONSTRAINT "message_text_next_fkey" FOREIGN KEY (next) REFERENCES message_text(id) ON UPDATE CASCADE ON DELETE CASCADE NOT DEFERRABLE;