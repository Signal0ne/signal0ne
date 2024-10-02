import { components, GroupBase, SingleValueProps } from 'react-select';
import { IncidentAssignee } from '../../../../contexts/IncidentsProvider/IncidentsProvider';
import { UserIcon } from '../../../Icons/Icons';
import './AssigneeDropdownSingleValue.scss';

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

const AssigneeDropdownSingleValue = (
  props: AssigneeDropdownSingleValueProps
) => (
  <components.SingleValue {...props}>
    {props.data?.value.photoUri ? (
      <img
        alt={`User: ${props.data?.value?.name}`}
        className="incident-task-assignee-avatar"
        src={props.data.value.photoUri}
        title={`User: ${props.data?.value?.name}`}
      />
    ) : (
      <UserIcon
        className="incident-task-assignee-avatar--empty"
        height={16}
        title={`User: ${props.data?.value?.name}`}
        width={16}
      />
    )}
  </components.SingleValue>
);

export default AssigneeDropdownSingleValue;
