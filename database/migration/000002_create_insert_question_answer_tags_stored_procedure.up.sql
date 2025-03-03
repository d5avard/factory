CREATE OR REPLACE FUNCTION insert_question_answer_tags(
    p_question_text TEXT,
    p_answer_text TEXT,
    p_tags TEXT[]
) RETURNS INT AS $$
DECLARE
    question_id INT;
    tag_id INT;
    tag TEXT;
BEGIN
    -- Insert Question and Get ID
    INSERT INTO questions (question_text) 
    VALUES (p_question_text)
    RETURNING id INTO question_id;

    -- Insert Answer
    INSERT INTO answers (question_id, answer_text) 
    VALUES (question_id, p_answer_text);

    -- Loop through Tags Array
    FOREACH tag IN ARRAY p_tags LOOP
        -- Insert Tag if Not Exists
        INSERT INTO tags (tag_name)
        VALUES (tag)
        ON CONFLICT (tag_name) DO NOTHING;

        -- Get Tag ID
        SELECT id INTO tag_id FROM tags WHERE tag_name = tag;

        -- Link Question with Tag
        INSERT INTO question_tags (question_id, tag_id)
        VALUES (question_id, tag_id)
        ON CONFLICT DO NOTHING;
    END LOOP;

    RETURN question_id;
END
$$ LANGUAGE plpgsql;

