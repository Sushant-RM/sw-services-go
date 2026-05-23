import json
import os

def generate_inserts():
    json_path = "scratch/wf_newsw1.json"
    if not os.path.exists(json_path):
        print(f"Error: {json_path} does not exist!")
        return

    with open(json_path, "r") as f:
        data = json.load(f)

    business_services = data.get("BusinessServices", [])
    if not business_services:
        print("Error: No BusinessServices found in JSON!")
        return

    sql_statements = []
    
    # We will register for both tenants: 'pb.amritsar' and 'pb'
    tenants = ['pb.amritsar', 'pb']
    
    for t_id in tenants:
        sql_statements.append(f"DELETE FROM eg_wf_action_v2 WHERE tenantid = '{t_id}';")

    for bs in business_services:
        orig_bs_uuid = bs.get("uuid")
        business_service = bs.get("businessService")
        business = bs.get("business")
        sla = bs.get("businessServiceSla")
        
        bs_audit = bs.get("auditDetails", {})
        created_by = bs_audit.get("createdBy", "sys-admin-uuid")
        created_time = bs_audit.get("createdTime", 1779517041553)
        modified_by = bs_audit.get("lastModifiedBy", "sys-admin-uuid")
        modified_time = bs_audit.get("lastModifiedTime", 1779517041553)

        for t_id in tenants:
            # Avoid UUID conflicts across tenants by appending suffix if not original tenant
            bs_uuid = orig_bs_uuid if t_id == 'pb.amritsar' else f"{orig_bs_uuid}_pb"
            
            sql_statements.append(f"DELETE FROM eg_wf_state_v2 WHERE tenantid = '{t_id}' AND businessserviceid = '{bs_uuid}';")
            sql_statements.append(f"DELETE FROM eg_wf_businessservice_v2 WHERE tenantid = '{t_id}' AND businessservice = '{business_service}';")

            # Insert business service
            sql_statements.append(
                f"INSERT INTO eg_wf_businessservice_v2 (uuid, tenantid, businessservice, business, createdby, createdtime, lastmodifiedby, lastmodifiedtime, businessservicesla) "
                f"VALUES ('{bs_uuid}', '{t_id}', '{business_service}', '{business}', '{created_by}', {created_time}, '{modified_by}', {modified_time}, {sla if sla is not None else 'NULL'});"
            )

            states = bs.get("states", [])
            
            # Map state name to tenant-specific UUID
            state_name_to_uuid = {}
            for state in states:
                s_name = state.get("state")
                orig_s_uuid = state.get("uuid")
                s_uuid = orig_s_uuid if t_id == 'pb.amritsar' else f"{orig_s_uuid}_pb"
                if s_name:
                    state_name_to_uuid[s_name] = s_uuid
                else:
                    state_name_to_uuid["START"] = s_uuid

            seq = 0
            for state in states:
                orig_s_uuid = state.get("uuid")
                s_uuid = orig_s_uuid if t_id == 'pb.amritsar' else f"{orig_s_uuid}_pb"
                s_name = state.get("state")
                app_status = state.get("applicationStatus")
                s_sla = state.get("sla")
                doc_req = state.get("docUploadRequired", False)
                start_state = state.get("isStartState", False)
                term_state = state.get("isTerminateState", False)
                updatable = state.get("isStateUpdatable", False)
                
                s_audit = state.get("auditDetails", {})
                s_created_by = s_audit.get("createdBy", created_by)
                s_created_time = s_audit.get("createdTime", created_time)
                s_modified_by = s_audit.get("lastModifiedBy", modified_by)
                s_modified_time = s_audit.get("lastModifiedTime", modified_time)

                # Insert state
                sql_statements.append(
                    f"INSERT INTO eg_wf_state_v2 (uuid, tenantid, businessserviceid, state, applicationstatus, sla, docuploadrequired, isstartstate, isterminatestate, createdby, createdtime, lastmodifiedby, lastmodifiedtime, seq, isstateupdatable) "
                    f"VALUES ('{s_uuid}', '{t_id}', '{bs_uuid}', "
                    f"'{s_name}' if '{s_name}' != 'None' else NULL, "
                    f"'{app_status}' if '{app_status}' != 'None' else NULL, "
                    f"{s_sla if s_sla is not None else 'NULL'}, "
                    f"{doc_req}, {start_state}, {term_state}, "
                    f"'{s_created_by}', {s_created_time}, '{s_modified_by}', {s_modified_time}, {seq}, {updatable});"
                )
                seq += 1

                actions = state.get("actions", [])
                if actions:
                    for act in actions:
                        orig_act_uuid = act.get("uuid")
                        act_uuid = orig_act_uuid if t_id == 'pb.amritsar' else f"{orig_act_uuid}_pb"
                        action_name = act.get("action")
                        next_state_val = act.get("nextState")
                        
                        # Resolve next state UUID
                        next_state_uuid = state_name_to_uuid.get(next_state_val, next_state_val)
                        if next_state_uuid and t_id == 'pb' and not next_state_uuid.endswith("_pb") and next_state_uuid != "None":
                            next_state_uuid = f"{next_state_uuid}_pb"
                        
                        roles_list = act.get("roles", [])
                        if not isinstance(roles_list, list):
                            roles_list = [roles_list] if roles_list else []
                        roles_list = list(roles_list)
                        if "SUPERUSER" not in roles_list:
                            roles_list.append("SUPERUSER")
                        roles = ",".join(roles_list)
                        active = act.get("active", True)
                        
                        act_audit = act.get("auditDetails", {})
                        a_created_by = act_audit.get("createdBy", created_by)
                        a_created_time = act_audit.get("createdTime", created_time)
                        a_modified_by = act_audit.get("lastModifiedBy", modified_by)
                        a_modified_time = act_audit.get("lastModifiedTime", modified_time)

                        # Insert action
                        sql_statements.append(
                            f"INSERT INTO eg_wf_action_v2 (uuid, tenantid, currentstate, action, nextstate, roles, createdby, createdtime, lastmodifiedby, lastmodifiedtime, active) "
                            f"VALUES ('{act_uuid}', '{t_id}', '{s_uuid}', '{action_name}', '{next_state_uuid}', '{roles}', '{a_created_by}', {a_created_time}, '{a_modified_by}', {a_modified_time}, {active});"
                        )

    # Clean up SQL formatting for None/NULL strings
    cleaned_sql = []
    for stmt in sql_statements:
        stmt = stmt.replace("'None' if 'None' != 'None' else NULL", "NULL")
        stmt = stmt.replace("True", "true").replace("False", "false")
        if " if " in stmt and " else " in stmt:
            parts = stmt.split(",")
            new_parts = []
            for p in parts:
                if " if " in p and " else " in p:
                    subparts = p.strip().split(" if ")
                    val = subparts[0]
                    cond_val = subparts[1].split(" != ")[0]
                    if cond_val == "'None'":
                        new_parts.append("NULL")
                    else:
                        new_parts.append(val)
                else:
                    new_parts.append(p)
            stmt = ", ".join(new_parts)
        cleaned_sql.append(stmt)

    with open("scratch/wf_insert.sql", "w") as f:
        f.write("\n".join(cleaned_sql))
    
    print("SQL file generated successfully with double-tenant residency at scratch/wf_insert.sql")

if __name__ == "__main__":
    generate_inserts()
