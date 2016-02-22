CREATE OR REPLACE FUNCTION insert_or_update_relationship(
	arg_id BIGINT,
	arg_other_id BIGINT,
	arg_state VARCHAR(8)
) RETURNS TABLE (
    id BIGINT,
    other_id BIGINT,
    state VARCHAR(8),
    type VARCHAR(16)
) AS $$
DECLARE
    query_type VARCHAR(16) = 'relationship';
	my_prev_state VARCHAR(8);	    --我方之前的状态
	other_prev_state VARCHAR(8);    --对方之前的状态
	final_state VARCHAR(8);		    --最终状态（用作函数返回值）
BEGIN
    SELECT INTO my_prev_state R.state FROM relationships AS R WHERE R.id=arg_id AND R.other_id=arg_other_id;
    SELECT INTO other_prev_state R.state FROM relationships AS R WHERE R.id=arg_other_id AND R.other_id=arg_id;
	IF arg_state = 'liked' THEN
		IF my_prev_state is null THEN	--我方条目尚未存在
			IF other_prev_state = 'liked' THEN	--对方条目已经存在并且为liked
				INSERT INTO relationships VALUES (arg_id, arg_other_id, 'matched');
				UPDATE relationships AS R SET state='matched' WHERE R.id=arg_other_id AND R.other_id=arg_id;
				final_state = 'matched';
			ELSE	--对方条目尚未存在或不是liked
				INSERT INTO relationships VALUES (arg_id, arg_other_id, 'liked');
				final_state = 'liked';
			END IF;
		ELSE	--我方条目已经存在
			IF other_prev_state = 'liked' THEN	--对方条目已经存在并且为liked
				UPDATE relationships AS R SET state='matched' WHERE R.id=arg_id AND R.other_id=arg_other_id;
				UPDATE relationships AS R SET state='matched' WHERE R.id=arg_other_id AND R.other_id=arg_id; --对方条目由liked升级为matched
				final_state = 'matched';
			ELSE 	--对方条目尚未存在或不是liked
				IF my_prev_state <> 'liked' THEN	-- 我方条目已经存在，并且不是liked
					UPDATE relationships AS R SET state='liked' WHERE R.id=arg_id AND R.other_id=arg_other_id;
				END IF;
				final_state = 'liked';
			END IF;
		END IF;
	ELSEIF arg_state = 'disliked' THEN
		IF my_prev_state is null THEN	--我方条目尚未存在
			IF other_prev_state = 'matched' THEN	--对方条目已经存在并且为matched
				INSERT INTO relationships VALUES (arg_id, arg_other_id, 'disliked');
				UPDATE relationships AS R SET state='liked' WHERE R.id=arg_other_id AND R.other_id=arg_id;
			ELSE	--对方条目尚未存在或不是matched
				INSERT INTO relationships VALUES (arg_id, arg_other_id, 'disliked');
			END IF;
		ELSE	--我方条目已经存在
			IF other_prev_state = 'matched' THEN	--对方条目已经存在并且为matched
				UPDATE relationships AS R SET state='disliked' WHERE R.id=arg_id AND R.other_id=arg_other_id;
				UPDATE relationships AS R SET state='liked' WHERE R.id=arg_other_id AND R.other_id=arg_id;	--对方条目由matched退化为liked
			ELSE 	--对方条目尚未存在或不是matched
				IF my_prev_state <> 'disliked' THEN 	--我方条目已经存在，并且不是disliked
					UPDATE relationships AS R SET state='disliked' WHERE R.id=arg_id AND R.other_id=arg_other_id;
				END IF;
			END IF;
		END IF;
		final_state = 'disliked';
	ELSE
		--APP SHOULD NOT MAKE THIS HAPPEN
	END IF;
	RETURN QUERY SELECT arg_id, arg_other_id, final_state, query_type;
END; $$
LANGUAGE plpgsql;
