ALTER TABLE management.label_policies ALTER COLUMN label_policy_state SET DEFAULT 0;
ALTER TABLE management.label_policies ALTER COLUMN label_policy_state SET NOT NULL;
ALTER TABLE management.label_policies ADD COLUMN warn_color STRING;
ALTER TABLE management.label_policies ADD COLUMN primary_color_dark STRING;
ALTER TABLE management.label_policies ADD COLUMN secondary_color_dark STRING;
ALTER TABLE management.label_policies ADD COLUMN warn_color_dark STRING;
ALTER TABLE management.label_policies ADD COLUMN logo_url STRING;
ALTER TABLE management.label_policies ADD COLUMN icon_url STRING;
ALTER TABLE management.label_policies ADD COLUMN logo_dark_url STRING;
ALTER TABLE management.label_policies ADD COLUMN icon_dark_url STRING;
ALTER TABLE management.label_policies ADD COLUMN err_msg_popup BOOLEAN;
ALTER TABLE management.label_policies ADD COLUMN disable_watermark BOOLEAN;

ALTER TABLE adminapi.label_policies ALTER COLUMN label_policy_state SET DEFAULT 0;
ALTER TABLE adminapi.label_policies ALTER COLUMN label_policy_state SET NOT NULL;
ALTER TABLE adminapi.label_policies ADD COLUMN warn_color STRING;
ALTER TABLE adminapi.label_policies ADD COLUMN primary_color_dark STRING;
ALTER TABLE adminapi.label_policies ADD COLUMN secondary_color_dark STRING;
ALTER TABLE adminapi.label_policies ADD COLUMN warn_color_dark STRING;
ALTER TABLE adminapi.label_policies ADD COLUMN logo_url STRING;
ALTER TABLE adminapi.label_policies ADD COLUMN icon_url STRING;
ALTER TABLE adminapi.label_policies ADD COLUMN logo_dark_url STRING;
ALTER TABLE adminapi.label_policies ADD COLUMN icon_dark_url STRING;
ALTER TABLE adminapi.label_policies ADD COLUMN err_msg_popup BOOLEAN;
ALTER TABLE adminapi.label_policies ADD COLUMN disable_watermark BOOLEAN;

ALTER TABLE auth.label_policies ALTER COLUMN label_policy_state SET DEFAULT 0;
ALTER TABLE auth.label_policies ALTER COLUMN label_policy_state SET NOT NULL;
ALTER TABLE auth.label_policies ADD COLUMN warn_color STRING;
ALTER TABLE auth.label_policies ADD COLUMN primary_color_dark STRING;
ALTER TABLE auth.label_policies ADD COLUMN secondary_color_dark STRING;
ALTER TABLE auth.label_policies ADD COLUMN warn_color_dark STRING;
ALTER TABLE auth.label_policies ADD COLUMN logo_url STRING;
ALTER TABLE auth.label_policies ADD COLUMN icon_url STRING;
ALTER TABLE auth.label_policies ADD COLUMN logo_dark_url STRING;
ALTER TABLE auth.label_policies ADD COLUMN icon_dark_url STRING;
ALTER TABLE auth.label_policies ADD COLUMN err_msg_popup BOOLEAN;
ALTER TABLE auth.label_policies ADD COLUMN disable_watermark BOOLEAN;


BEGIN;
ALTER TABLE management.label_policies DROP CONSTRAINT "primary";
ALTER TABLE management.label_policies ADD CONSTRAINT "primary" PRIMARY KEY (aggregate_id, label_policy_state);
ALTER TABLE adminapi.label_policies DROP CONSTRAINT "primary";
ALTER TABLE adminapi.label_policies ADD CONSTRAINT "primary" PRIMARY KEY (aggregate_id, label_policy_state);
ALTER TABLE auth.label_policies DROP CONSTRAINT "primary";
ALTER TABLE auth.label_policies ADD CONSTRAINT "primary" PRIMARY KEY (aggregate_id, label_policy_state);
COMMIT;


ALTER TABLE management.users ADD COLUMN avatar STRING;
ALTER TABLE auth.users ADD COLUMN avatar STRING;
ALTER TABLE adminapi.users ADD COLUMN avatar STRING;