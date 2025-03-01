import {
  Button,
  ButtonGroup,
  Container,
  Form,
  Modal,
  Row,
} from "react-bootstrap";
import { CommandType, HeaderType, NotifyType, WebHookType } from "types/config";
import { FormProvider, useForm } from "react-hook-form";
import { flattenErrors, removeEmptyValues } from "utils";
import { useCallback, useContext, useState } from "react";

import { DeleteModal } from "./delete-confirm";
import EditService from "components/modals/service-edit/service";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { HelpTooltip } from "components/generic";
import { ModalContext } from "contexts/modal";
import { ServiceEditType } from "types/service-edit";
import { convertUIServiceDataEditToAPI } from "components/modals/service-edit/util/ui-api-conversions";
import { faCircleNotch } from "@fortawesome/free-solid-svg-icons";
import { useTheme } from "contexts/theme";

export interface EditForm {
  optionsSemanticVersioning?: boolean;
  lvType?: "github" | "url";
  lvAllowInvalidCerts?: boolean;
  lvUsePreReleases?: boolean;
  dvCustomHeaders?: HeaderType[];
  dvAllowInvalidCerts?: boolean;
  commands?: CommandType[];
  webhooks?: WebHookType[];
  notify?: NotifyType[];
  dashboardAutoApprove?: boolean;
}

const getPayload = (data: ServiceEditType) => {
  return removeEmptyValues(convertUIServiceDataEditToAPI(data));
};

const ServiceEditModal = () => {
  // modal.actionType:
  // EDIT
  const { handleModal, modal } = useContext(ModalContext);
  const themeCtx = useTheme();
  const form = useForm<ServiceEditType>({ mode: "onBlur" });
  // null if submitting
  const [err, setErr] = useState<string | null>("");

  const hideModal = useCallback(() => {
    form.reset({});
    setErr("");
    handleModal("", { id: "", loading: true });
  }, []);

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const onSubmit = async (data: ServiceEditType) => {
    setErr(null);
    const payload = getPayload(data);
    const serviceName = modal.service.id;

    await fetch(
      serviceName
        ? `/api/v1/service/edit/${serviceName}`
        : "/api/v1/service/new",
      {
        method: serviceName ? "PUT" : "POST",
        body: JSON.stringify(payload),
      }
    )
      .then((response) => {
        if (!response.ok) throw response;
        hideModal();
      })
      .catch(async (err) => {
        let errorMessage = err.statusText;
        try {
          const responseBody = await err.json();
          errorMessage = responseBody.message;
          setErr(errorMessage);
        } catch (e) {
          console.log(e);
          setErr(err.toString());
        }
      });
  };

  const onDelete = async () => {
    console.log(`Deleting ${modal.service.id}`);
    await fetch(`/api/v1/service/delete/${modal.service.id}`, {
      method: "DELETE",
    }).then(() => {
      hideModal();
    });
  };

  return (
    <FormProvider {...form}>
      <Form id="service-edit">
        <Modal
          size="lg"
          show={modal.actionType === "EDIT"}
          onHide={() => hideModal()}
        >
          <Modal.Header
            closeButton
            closeVariant={themeCtx.theme === "theme-dark" ? "white" : undefined}
          >
            <Modal.Title>
              <strong>Edit Service</strong>
              <HelpTooltip
                text="Greyed out placeholder text represents a default that you can override. (current secrets can be kept by leaving them as '<secret>')"
                placement="bottom"
              />
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <Container
              fluid
              className="font-weight-bold"
              style={{ paddingLeft: "0px" }}
            >
              <EditService name={modal.service.id} />
            </Container>
          </Modal.Body>
          <Modal.Footer
            style={{ display: "flex", justifyContent: "space-between" }}
          >
            <ButtonGroup>
              <DeleteModal
                onDelete={() => onDelete()}
                disabled={err === null}
              />
            </ButtonGroup>
            {err === null && (
              <FontAwesomeIcon
                icon={faCircleNotch}
                style={{
                  padding: "0",
                }}
                className="fa-spin"
              />
            )}
            <span>
              <Button
                id="modal-cancel"
                variant="secondary"
                onClick={() => hideModal()}
                disabled={err === null}
              >
                Cancel
              </Button>
              <Button
                id="modal-action"
                variant="primary"
                type="submit"
                onClick={form.handleSubmit(onSubmit)}
                className="ms-2"
                disabled={err === null}
              >
                Confirm
              </Button>
            </span>
            {form.formState.submitCount > 0 &&
              (!form.formState.isValid || err) && (
                <Row>
                  <div className="error-msg">
                    Please correct the errors in the form and try again.
                    <br />
                    {/* Render either the server error or form validation error */}
                    {err ? (
                      <>
                        <>
                          {err.split("\\").map((line) => (
                            <pre key={line} className="no-margin">
                              {line}
                            </pre>
                          ))}
                        </>
                      </>
                    ) : (
                      <ul>
                        {Object.entries(
                          flattenErrors(form.formState.errors)
                        ).map(([key, error]) => (
                          <li key={key}>
                            {key}: {error}
                          </li>
                        ))}
                      </ul>
                    )}
                  </div>
                </Row>
              )}
          </Modal.Footer>
        </Modal>
      </Form>
    </FormProvider>
  );
};

export default ServiceEditModal;
