package secure_dorm

/*
 * Capability is a structure that encodes information about which
 * objects a user can read or write.
 */
type Capability struct {
}

/*
 * Given a capability and an object, calling cap.CanRead(object) returns
 * true if the capability permits the user to read the object. CanRead
 * expects that its argument is a pointer to a struct.
 *
 * To be explicit, `object` will have type: *MyStruct,
 * where MyStruct is any arbitrary struct subject to the restrictions
 * discussed later in this document.
 *
 * Example usage to test if the caller can read a post object:
 *    type Post struct = { ... }
 *    cap := ...
 *    post := &Post{}
 *    ok := cap.CanRead(post)
 */
func (cap *Capability) CanRead(object interface{}) bool {
	return false
}

/*
 * Given a capability and an object, calling cap.CanWrite(object) returns
 * true if the capability permits the user to write to the object.
 *
 * As mentioned in the description of CanRead, the argument `object`
 * will be a pointer to a model.
 */
func (cap *Capability) CanWrite(object interface{}) bool {
	return false
}

/*
 * The capability manager allows users to create and modify capabilities.
 */
type CapabilityManager struct {
}

/*
 * Creates a new instance of a capability manager.
 */
func NewCapabilityManager() *CapabilityManager {
	return &CapabilityManager{}
}

/*
 * Given a unique username, cm.GetRootCapability(username) returns the user's
 * root capability (or nil if one has not yet been set). A root capability bootstraps
 * the user's permissions. For instance, a newly created user's root capability might
 * just include the ability to read and write their own user object.
 */
func (cm *CapabilityManager) GetRootCapability(username string) *Capability {
	return nil
}

/*
 * A root capability bootstraps a user's permissions. Given a unique username and
 * two slices of objects, cm.SetRootCapability(username, readSet, writeSet)
 * associates a root capability with the username. The root capability is expected to allow
 * reading and writing all objects in readSet and writeSet, respectively. For instance,
 * a newly created user's root capability might just include the ability to read and write
 * their own user object. Thus, after creating the new object, the user's root capability
 * would be set by `cm.SetRootCapability(user.Username, []interface{user}, []interface{user}).`
 */
func (cm *CapabilityManager) SetRootCapability(username string,
	readSet []interface{}, writeSet []interface{}) {
}

/*
 * Given a capability and an object, cm.AddReadCapability(cap, object) returns a new capability
 * that includes all capabilities of cap plus the ability to read object. That is, if newCap is
 * the new capability, then calling newCap.CanRead(object) should return true. Note, however, that
 * the original capability should not be modified, so calling cap.CanRead(object) should still
 * return false. Similarly, root capabilities should not change.
 */
func (cm *CapabilityManager) AddReadCapability(cap *Capability, object interface{}) *Capability {
	return nil
}

/*
 * Given a capability and an object, cm.AddWriteapCability(cap, object) returns a new capability
 * that includes all capabilities of cap plus the ability to write object. Like mentioned above
 * for `AddReadCapability()`, the original capability and all root capabilities should not be modified.
 */
func (cm *CapabilityManager) AddWriteCapability(cap *Capability, object interface{}) *Capability {
	return nil
}

/*
 * Given a capability and an object, cm.RemoveReadCapability(cap, object) returns a new capability
 * that includes all capabilities of cap minus the ability to read object. Like mentioned above
 * for `AddReadCapability()`, the original capability and all root capabilities should not be modified.
 */
func (cm *CapabilityManager) RemoveReadCapability(cap *Capability, object interface{}) *Capability {
	return nil
}

/*
 * Given a capability and an object, cm.RemoveWriteCapability(cap, object) returns a new capability
 * that includes all capabilities of cap minus the ability to write object. Like mentioned above
 * for `AddReadCapability()`, the original capability and all root capabilities should not be modified.
 */
func (cm *CapabilityManager) RemoveWriteCapability(cap *Capability, object interface{}) *Capability {
	return nil
}
