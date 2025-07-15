import { CellContext } from '@tanstack/react-table';
import { Edit, X } from 'lucide-react';

import { useRbacRoles } from '@/react/portainer/users/RolesView/useRbacRoles';

import { Button } from '@@/buttons';
import { Select } from '@@/form-components/Input';

import { Access, getTableMeta } from '../types';

import { helper } from './helper';
import { RoleService } from 'Portainer/rbac/services/role.service';

export const role = helper.accessor('Role', {
  cell: RoleCell,
  header: 'Role',
  meta: {
    width: 320,
  },
});

function RoleCell({
  row: { original: item, getCanSelect },
  table,
  getValue,
}: CellContext<Access, string>) {
  const meta = getTableMeta(table.options.meta);
  const type = item.Type as 'team' | 'user';
  const updateValue = meta.roles.getRoleValue(item.Id, type);
  const roleId = Number(getValue());

  if (!getCanSelect()) {
    return <>{roleId}</>;
  }

  const roles = RoleService().roles();
  const role = roles.find((r) => r.ID === roleId);
  const roleName = role ? role.Name : 'Unknown';

  if (typeof updateValue === 'undefined') {
    return (
      <>
        { roleName }
        <Button
          color="none"
          icon={Edit}
          onClick={() => meta.roles.setRolesValue(item.Id, type, roleId)}
          data-cy="edit-role-button"
        >
          Edit
        </Button>
      </>
    );
  }
  return (
    <RollEdit
      value={updateValue}
      onChange={(value) => meta.roles.setRolesValue(item.Id, type, value)}
    />
  );
}

function RollEdit({
  value,
  onChange,
}: {
  value: number;
  onChange(value?: number): void;
}) {
  const rolesQuery = useRbacRoles({
    select: (roles) => roles.map((r) => ({ label: r.Name, value: r.Id })),
  });

  if (!rolesQuery.data) {
    return null;
  }

  return (
    <div className="flex items-center gap-3 max-w-xs">
      <Select
        aria-label="Role"
        data-cy="role-select"
        value={value}
        options={rolesQuery.data}
        onChange={(e) => onChange(parseInt(e.target.value, 10))}
      />
      <Button
        color="none"
        icon={X}
        onClick={() => onChange()}
        data-cy="cancel-role-button"
      />
    </div>
  );
}
