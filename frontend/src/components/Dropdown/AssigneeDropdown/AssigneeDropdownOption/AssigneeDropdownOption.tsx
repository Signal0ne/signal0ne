import { components, GroupBase, OptionProps } from 'react-select';
import { IncidentAssignee } from '../../../../contexts/IncidentsProvider/IncidentsProvider';
import './AssigneeDropdownOption.scss';

type AssigneeDropdownOptionProps = OptionProps<
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

const AssigneeDropdownOption = ({
  children,
  ...props
}: AssigneeDropdownOptionProps) => (
  <components.Option {...props}>
    <img
      alt={`User: ${props.data?.value?.name}`}
      className="incident-task-assignee-avatar"
      src={props.data?.value?.photoUri}
    />
    <span className="incident-task-assignee-name">{children}</span>
  </components.Option>
);

export default AssigneeDropdownOption;
