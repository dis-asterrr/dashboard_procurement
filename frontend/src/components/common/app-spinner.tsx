"use client";

import { Card, Flex, Spin, Typography, theme } from "antd";
import { LoadingOutlined } from "@ant-design/icons";

type AppSpinnerProps = {
  text?: string;
  card?: boolean;
  size?: number;
};

const { Text } = Typography;

export default function AppSpinner({
  text = "Loading...",
  card = false,
  size = 30,
}: AppSpinnerProps) {
  const { token } = theme.useToken();
  const indicator = <LoadingOutlined style={{ fontSize: size, color: token.colorPrimary }} spin />;

  if (card) {
    return (
      <Card
        variant="borderless"
        style={{
          minWidth: 280,
          borderRadius: 12,
          border: `1px solid ${token.colorBorderSecondary}`,
          boxShadow: token.boxShadowTertiary,
          background: token.colorBgElevated,
        }}
        styles={{ body: { padding: "24px 28px" } }}
      >
        <Flex vertical align="center" gap={10}>
          <Spin indicator={indicator} />
          <Text type="secondary" style={{ fontSize: 14 }}>
            {text}
          </Text>
        </Flex>
      </Card>
    );
  }

  return (
    <Flex vertical align="center" justify="center" gap={10}>
      <Spin indicator={indicator} />
      {text ? (
        <Text type="secondary" style={{ fontSize: 14 }}>
          {text}
        </Text>
      ) : null}
    </Flex>
  );
}
