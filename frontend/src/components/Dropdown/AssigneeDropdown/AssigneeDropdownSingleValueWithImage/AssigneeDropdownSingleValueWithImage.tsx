import { components, GroupBase, SingleValueProps } from 'react-select';
import { IncidentAssignee } from '../../../../contexts/IncidentsProvider/IncidentsProvider';
import './AssigneeDropdownSingleValueWithImage.scss';

type AssigneeDropdownSingleValueProps = SingleValueProps<
  {
    label: string;
    value: IncidentAssignee;
  },
  false,
  GroupBase<{
    label: string;
    value: IncidentAssignee;
  }>
>;

const AssigneeDropdownSingleValueWithImage = (
  props: AssigneeDropdownSingleValueProps
) => (
  <components.SingleValue {...props}>
    {props.data?.value?.photoUri && (
      <img
        alt={`User: ${props.data?.value?.name}`}
        className="incident-task-assignee-avatar-image"
        src={props.data.value.photoUri}
        title={`User: ${props.data?.value?.name}`}
      />
    )}
    <span>{props.data?.value?.name}</span>
  </components.SingleValue>
);

export default AssigneeDropdownSingleValueWithImage;
