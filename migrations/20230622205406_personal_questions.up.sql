CREATE TABLE IF NOT EXISTS personal_questions (
	uuid uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
	/*----------------------------*/
    -- extras
    question varchar
);

INSERT INTO personal_questions(uuid, question)
values ('288d830d-b3cc-4320-995b-27d3ecf99b24','To what city did you go the first time you flew on a plane ?'),
       ('2bce22f5-dd6f-4536-90f6-13f46f2cf94b','What is the first name of the person you first kissed ?'),
       ('3573d05f-ee65-409e-aaf2-03495dbfeb98','What is the first name of your best friend in high school ?'),
       ('503d834c-b44c-4780-801b-9a58a3cae8d8','What is the first name of your oldest nephew ?'),
       ('573cd9e6-e3e0-46af-9e10-0d9ed6d38bd5','What is the first name of your oldest niece ?'),
       ('5f1b9875-38df-4b12-b1e9-cec9d0c39f7b','What is the first name of your oldest niece ?'),
       ('5fded606-c632-4030-9d85-bfc146d16d21','What is your oldest cousin''s first name ?'),
       ('93e324d7-c27d-4a0f-bdf7-bb75583a0ae8','What was the first name of the first person you dated ?'),
       ('9b0050d4-55b7-44e7-8e27-6529a4acc9bd','What was the first name of your favorite childhood friend ?'),
       ('b402cf88-ae8c-4250-a65f-3050fdbfc439','What was the last name of your third grade teacher ?'),
       ('c710d94e-6da9-4d45-a8d3-087bdc264adb','What was the street name where your best friend in high school lived (street name only) ?'),
       ('ded4f973-7a95-4aab-b5fc-4093dab906ab','Where were you when you first heard about 9/11 ?'),
       ('faa070d2-4ff2-4b77-aa3a-3f15b513bd0e','What was the first name of the boy/girl with whom you had your second kiss ?')

