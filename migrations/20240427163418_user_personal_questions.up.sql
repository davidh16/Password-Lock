CREATE TABLE IF NOT EXISTS user_personal_questions (
    user_uuid uuid references users(uuid),
    personal_question_uuid uuid references personal_questions(uuid),
    answer varchar
);
