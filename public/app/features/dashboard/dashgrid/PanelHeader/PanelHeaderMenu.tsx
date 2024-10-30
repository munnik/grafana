import * as React from 'react';

import { PanelMenuItem } from '@grafana/data';
import { selectors } from '@grafana/e2e-selectors';
import { Menu } from '@grafana/ui';

export interface Props {
  items: PanelMenuItem[];
  style?: React.CSSProperties;
  itemsClassName?: string;
  className?: string;
}

export function PanelHeaderMenu({ items }: Props) {
  return <Menu>{renderPanelMenuItems(items)}</Menu>;
}

export const renderPanelMenuItems = (items: PanelMenuItem[]) => {
  return items.map((item) => {
    if (item.component) {
      return <item.component key={item.text} />;
    }

    switch (item.type) {
      case 'divider':
        return <Menu.Divider key={item.text} />;
      case 'group':
        return (
          <Menu.Group key={item.text} label={item.text}>
            {item.subMenu ? renderPanelMenuItems(item.subMenu) : undefined}
          </Menu.Group>
        );
      default:
        return (
          <Menu.Item
            key={item.text}
            label={item.text}
            icon={item.iconClassName}
            childItems={item.subMenu ? renderPanelMenuItems(item.subMenu) : undefined}
            url={item.href}
            onClick={item.onClick}
            shortcut={item.shortcut}
            testId={selectors.components.Panels.Panel.menuItems(item.text)}
          />
        );
    }
  });
};
