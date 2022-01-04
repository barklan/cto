import { acceptHMRUpdate, defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  /**
   * Current named of the user.
   */
  const name = ref('')
  const id = ref('')
  const project = ref('')
  const token = ref('')

  /**
   *
   * @param name - new name to set
   */
  function setName(newName: string) {
    name.value = newName
  }

  function setID(newID: string) {
    id.value = newID
  }
  function setProject(newProject: string) {
    project.value = newProject
  }
  function setToken(newToken: string) {
    token.value = newToken
  }

  return {
    setName,
    setID,
    setProject,
    setToken,
    name,
    id,
    project,
    token,
  }
})

if (import.meta.hot)
  import.meta.hot.accept(acceptHMRUpdate(useUserStore, import.meta.hot))
