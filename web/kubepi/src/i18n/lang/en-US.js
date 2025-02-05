import el from "element-ui/lib/locale/lang/en"
import fu from "fit2cloud-ui/src/locale/lang/en_US" // 加载fit2cloud的内容

const message = {
  commons: {
    message_box: {
      alert: "Alter",
      confirm: "Confirm",
      prompt: "Prompt",
    },
    personal: {
      profile: "Profile",
      exit: "Logout",
      project_url: "Project Address",
      issue: "Feedback",
      talk: "Participate in the discussion",
      star: "Star",
      version: "Version",
      introduction: "KubePi is an easy to use open source Kubernetes visual management panel.",
      introduction2: "KubePi allows administrators to import multiple Kubernetes clusters and assign permissions of different clusters and namespaces to specific users through permission control. It allows developers to manage and troubleshoot applications running in the Kubernetes cluster to better deal with the complexity in the Kubernetes cluster."
    },
    help: {
      help: "Help",
      about: "About",
      business_support: "Business Support",
      ko_docs: "Official Document",
    },
    button: {
      detail: "Detail",
      delete: "Delete",
      skip: "Skip",
      import: "Import",
      create: "Create",
      cancel: "Cancel",
      login: "Login",
      confirm: "Confirm",
      add: "Add",
      edit: "Edit",
      all_select: "All",
      upload: "Upload",
      search: "Search",
      rbac_manage: "RBAC Manage"
    },
    table: {
      name: "name",
      kind: "kind",
      created_time: "created at",
      status: "status",
      action: "action",
      creat_by: "created by",
      built_in: "built in",
      description: "description",
    },
    header: {
      help_doc: "document",
      support: "support",
      guide: "guide",
      guide_text: "Thank you for choosing this product. Would you like to register your first cluster now?"
    },
    bool: {
      true: "true",
      false: "false"
    },
    form: {
      select_placeholder: "please select"
    },
    validate: {
      limit: "长度在 {0} 到 {1} 个字符",
      input: "please input {0}",
      select: "please select {0}",
      required: "required",
      email: "please input a valid email",
      password_help: "Valid password: 8-30 digits, English letters + numbers + special characters (optional)",
      name_not_compliant: "The name does not conform to the naming convention!",
    },
    msg: {
      create_success: "create success",
      delete_success: "delete success",
      update_success: "update success",
      no_data: "no data",
    },
    confirm_message: {
      delete: "This operation cannot be undone. Do you want to continue?",
      create_success: "create success",
      save_success: "save success",
    },
    login: {
      username_or_email: "username or email",
      password: "password",
      title: "login KubePi",
      welcome: "Welcome back, please enter your user name and password to log in",
      expires: "The authentication information has expired. Please log in again."
    },
  },
  business: {
    cluster: {
      cluster: "Cluster",
      namespace: "Namespace",
      scope: "scope",
      version: "version",
      list: "Cluster",
      import: "import cluster",
      edit: "edit",
      nodes: "nodes",
      label: "label",
      description: "description",
      cluster_detail: "Detail",
      connect_setting: "Connect setting",
      connect_direction: "Connect direction",
      connect_forward: "Forward",
      connect_backward: "Backward",
      authenticate_setting: "Authenticate setting",
      certificate: "Certificate",
      authenticate_mode: "Authenticate mode",
      expect: "expect",
      management: "management",
      open_dashboard: "Console",
      cluster_version: "version",
      member: "Member",
      role: "Role",
      administrator: "administrator",
      viewer: "viewer",
      custom: "custom",
      rule: "rule",
      config_file: "kubeconfig file",
      config_content: "Config content",
      hidden_cluster: "hide inaccessible clusters",
      namespace_role_form_check_msg: "namespace or role list cannot be empty",
      api_group: "API Groups",
      resource: "resources",
      verb: "verbs",
      cluster_role_form_check_msg: "api groups,resources or verbs  list cannot be empty",
      user_not_in_cluster: "You are not a member of this cluster and cannot access the console of this cluster. Please contact the administrator to add you as a cluster member.",
      ready: "Ready",
      not_ready: "NotReady"
    },
    cluster_role: {
      none: "None",
      cluster_administrator: "Administrator",
      cluster_viewer: "Viewer",
    },
    user: {
      user_management: "RBAC",
      username: "username",
      nickname: "nickname",
      email: "email",
      user_list: "User list",
      role_list: "Role list",
      user: "User",
      role: "Role",
      template: "template",
      base_on_exists_role: "base on exists role",
      permission: "permission",
      permission_setting: "permission setting",
      password: "password",
      confirm_password: "confirm password",
      old_password: "old password",
      new_password: "new password",
      change_password: "change password",
      resource_name: "resource name",
      please_input_password: "please input password",
      please_input_password_agin: "please input password again",
      password_not_equal: "two passwords are inconsistent"

    }
  },
}

const description = {
  i18n_user_administrator: "Super administrator, with permissions for all objects.",
  i18n_user_manage_cluster: "The Cluster Administrator has all the permissions of the cluster object.",
  i18n_user_manage_rbac: "Role and user administrators have all the permissions of role and user objects.",
  i18n_user_manage_readonly: "Read only user with access to all objects only,",
  i18n_user_common_user: "Ordinary users only have access to cluster objects",
  i18n_manage_nodes: "The node administrator has read-only permissions on node objects",
  i18n_manage_crd: "The custom resource administrator has all the permissions of the CustomResourceDefinition objec",
  i18n_view_crd: "User defined resource administrator with read-only permission on CustomResourceDefinition object",
  i18n_user_manage_chart:"Chart warehouse administrator, has all rights to the Chart warehouse",

  i18n_manage_config: "The configuration administrator has all permissions on the current namespace configmap, secret, resourcequotes, limitranges, horizontalpodautoscalers and poddisruptionbudget objects",
  i18n_view_config: "Configure a read-only user with read-only permissions on the current namespace configmap, secret, resourcequotes, limitranges, horizontalpodautoscalers and poddisruptionbudget objects",
  i18n_namespace_owner: "The namespace owner has all permissions on all objects in the current namespace",
  i18n_namespace_viewer: "Namespace read-only user with read-only permission for all objects in the current namespace",
  i18n_view_workload: "Workload read-only user with read-only permissions for daemonset, statefulset, deployment, job, cronjob and pod in the current namespace",
  i18n_manage_workload: "Workload administrator, read-only user of workload, with all permissions of daemonset, statefulset, deployment, job, cronjob and pod in the current namespace",
  i18n_manage_storage: "Storage administrator who has all permissions on the persistentvolumeclaim object in the current namespace",
  i18n_view_storage: "Stores a read-only user with read-only permissions on the persistentvolumeclaim object in the current namespace",
  i18n_view_service_discovery: "The service found a read-only user with read-only permissions on service, endpoint, progress and networkpolicy objects in the current namespace",
  i18n_manage_service_discovery: "The service discovery administrator has all permissions on service, endpoint, ingress and networkpolicy objects in the current namespace",
  i18n_manage_rbac: "The service discovery administrator has all permissions on service, endpoint, ingress and networkpolicy objects in the current namespace",
  i18n_view_rbac: "Namespace access control read-only user with read-only permissions for role, rolebinding and serviceaccount objects in the current namespace",
  i18n_manage_appmarket: "Application market administrator, who has all rights to the application market"
}

const apiObjects = {
  users: "users",
  roles: "roles",
  clusters: "clusters",
}

const apiVerbs = {
  "update": "update",
  "delete": "delete",
  "get": "get",
  "list": "list",
  "create": "create"
}


export default {
  ...el,
  ...fu,
  ...message,
  ...apiObjects,
  ...apiVerbs,
  ...description
}
