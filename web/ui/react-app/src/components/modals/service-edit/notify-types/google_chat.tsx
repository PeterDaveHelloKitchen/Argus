import { FormLabel, FormTextArea } from "components/generic/form";

import { NotifyGoogleChatType } from "types/config";
import { NotifyOptions } from "./generic";
import { useGlobalOrDefault } from "./util";

const GOOGLE_CHAT = ({
  name,

  global,
  defaults,
  hard_defaults,
}: {
  name: string;

  global?: NotifyGoogleChatType;
  defaults?: NotifyGoogleChatType;
  hard_defaults?: NotifyGoogleChatType;
}) => (
  <>
    <NotifyOptions
      name={name}
      global={global?.options}
      defaults={defaults?.options}
      hard_defaults={hard_defaults?.options}
    />
    <>
      <FormLabel text="URL Fields" heading />
      <FormTextArea
        name={`${name}.url_fields.raw`}
        required
        col_sm={12}
        rows={2}
        label="Raw"
        tooltip="e.g. chat.googleapis.com/v1/spaces/foo/messages?key=bar&token=baz"
        placeholder={useGlobalOrDefault(
          global?.url_fields?.raw,
          defaults?.url_fields?.raw,
          hard_defaults?.url_fields?.raw
        )}
      />
    </>
  </>
);

export default GOOGLE_CHAT;
