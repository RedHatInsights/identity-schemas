model InternalInfo {
  auth_time: number;
  cross_access?: boolean;
  org_id: string;
}

model UserIdentity {
  account_number: string;
  org_id: string;
  type: "User";
  auth_type: "jwt-auth" | "basic-auth";
  internal: InternalInfo;
  user: {
    email: string;
    first_name: string;
    last_name: string;
    is_active: boolean;
    is_internal: boolean;
    is_org_admin: boolean;
    locale: string;
    user_id: string;
    username: string;
  };
}

model ServiceAccountIdentity {
  account_number: string;
  org_id: string;
  type: "ServiceAccount";
  auth_type: "jwt-auth";
  internal: InternalInfo;
  service_account: {
    client_id: string;
    username: string;
    scope?: string;
    user_id: string;
  };
}

model CertSystemIdentity {
  account_number: string;
  org_id: string;
  type: "System";
  auth_type: "cert-auth";
  internal: InternalInfo;
  system: {
    cert_type: string;
    cn: string;
  };
}

model UHCSystemIdentity {
  account_number: string;
  org_id: string;
  type: "System";
  auth_type: "uhc-auth";
  internal: InternalInfo;
  system: {
    cluster_id: string;
  };
}

model Entitlement {
  is_entitled: boolean;
  is_trial: boolean;
}

model Entitlements {
  ansible?: Entitlement;
  cost_management?: Entitlement;
  insights?: Entitlement;
  internal?: Entitlement;
  migrations?: Entitlement;
  openshift?: Entitlement;
  rhel?: Entitlement;
  rhoam?: Entitlement;
  rhods?: Entitlement;
  rhosak?: Entitlement;
  settings?: Entitlement;
  smart_management?: Entitlement;
  subscriptions?: Entitlement;
  user_preferences?: Entitlement;
}

model Identity = UserIdentity | ServiceAccountIdentity | CertSystemIdentity | UHCSystemIdentity;

model ThreeScaleIdentitySchema {
  identity: Identity;
  entitlements: Entitlements;
} 