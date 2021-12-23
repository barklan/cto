import { acceptHMRUpdate, defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  /**
   * Current named of the user.
   */
  const name = ref('')
  const id = ref('')
  const role = ref('')
  const project = ref('')
  const token = ref('')

  const previousNames = ref(new Set<string>())

  const usedNames = computed(() => Array.from(previousNames.value))
  const otherNames = computed(() => usedNames.value.filter(n => n !== name.value))

  /**
   * Changes the current name of the user and saves the one that was used
   * before.
   *
   * @param name - new name to set
   */
  function setName(newName: string) {
    if (name.value)
      previousNames.value.add(name.value)

    name.value = newName
  }

  function setID(newID: string) {
    id.value = newID
  }
  function setProject(newProject: string) {
    project.value = newProject
  }
  function setRole(newRole: string) {
    role.value = newRole
  }
  function setToken(newToken: string) {
    token.value = newToken
  }

  return {
    setName,
    setID,
    setProject,
    setRole,
    setToken,
    otherNames,
    name,
    id,
    role,
    project,
    token,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))
